import { useUserStore } from '../stores/userStore';
import { useSessionStore } from '../stores/sessionStore';
import { useMessageStore } from '../stores/messageStore';
import { useWsStore } from '../stores/wsStore';

let pendingCommand = null;  // 'connect', 'join', 'create', 'send', 'exit'
let pendingSend = null;    // { sessionId, content } for pending send
let pendingCreateTheme = null; // theme name for pending create
let pendingExitId = null;  // session id for pending exit

function resolveSessionId(response) {
  const trimmed = String(response).trim();
  const id = parseInt(trimmed, 10);
  if (isNaN(id)) return null;
  return id;
}

export function connectWebSocket() {
  const { nickname, headshot, setUserInfo, clearUser } = useUserStore.getState();
  const { addJoinedSession } = useSessionStore.getState();
  const { clearAll: clearMessages } = useMessageStore.getState();
  const { setSocket, setConnected, setConnecting, clearSocket } = useWsStore.getState();

  if (!nickname) return;

  setConnecting(true);
  pendingCommand = 'connect';

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  const url = `${protocol}//${host}/connect?nickname=${encodeURIComponent(nickname)}&headshot=${encodeURIComponent(headshot)}`;

  const socket = new WebSocket(url);

  socket.onopen = () => {
    setConnecting(false);
    setConnected(true);
    setSocket(socket);
  };

  socket.onmessage = (event) => {
    const raw = typeof event.data === 'string' ? event.data : String(event.data);
    const data = raw.trim();

    if (data === 'y') {
      if (pendingCommand === 'send' && pendingSend) {
        const { sessionId, content } = pendingSend;
        const { uid, nickname, headshot } = useUserStore.getState();
        useMessageStore.getState().addMessage(sessionId, {
          uid,
          nickname,
          headshot,
          content,
          time: Date.now(),
          isSelf: true,
        });
        pendingSend = null;
      }
      if (pendingCommand === 'exit' && pendingExitId != null) {
        useSessionStore.getState().removeJoinedSession(pendingExitId);
        pendingExitId = null;
      }
      pendingCommand = null;
      return;
    }

    // Check if it's an incoming message: session_id\x1fuid\x1fmsg
    if (data.includes('\x1f')) {
      const parts = data.split('\x1f');
      if (parts.length >= 3) {
        const sessionId = parseInt(parts[0], 10);
        const senderUid = parseInt(parts[1], 10);
        const msg = parts.slice(2).join('\x1f');
        handleIncomingMessage(sessionId, senderUid, msg);
      }
      pendingCommand = null;
      return;
    }

    // Numeric response — could be uid or session_id
    const id = resolveSessionId(data);
    if (id !== null) {
      if (pendingCommand === 'connect') {
        setUserInfo(id, nickname, headshot);
        pendingCommand = null;
      } else if (pendingCommand === 'create') {
        const theme = pendingCreateTheme || null;
        addJoinedSession({ id, theme });
        pendingCommand = null;
        pendingCreateTheme = null;
      } else if (pendingCommand === 'join') {
        addJoinedSession({ id, theme: null });
        pendingCommand = null;
        fetchSessionTheme(id);
      }
    }
  };

  socket.onclose = () => {
    setConnected(false);
    setConnecting(false);
    clearSocket();
    pendingCommand = null;
    pendingSend = null;
    pendingCreateTheme = null;
    pendingExitId = null;

    // Full state reset on disconnect per PRD #5
    clearUser();
    clearMessages();
    useSessionStore.getState().clearAll();
  };

  socket.onerror = () => {
    // onclose will fire after this
  };

  return socket;
}

function handleIncomingMessage(sessionId, senderUid, content) {
  const { uid } = useUserStore.getState();
  const {
    addJoinedSession,
    updateSessionLastMessageTime,
    sortJoinedSessionsByRecentMessage,
    incrementUnread,
  } = useSessionStore.getState();
  const { addMessage, ensureSession } = useMessageStore.getState();

  const now = Date.now();

  // Ensure session exists in our state
  ensureSession(sessionId);
  const alreadyJoined = useSessionStore.getState().joinedSessions.find(
    (s) => s.id === sessionId
  );
  if (!alreadyJoined) {
    addJoinedSession({ id: sessionId, theme: null });
    fetchSessionTheme(sessionId);
  }

  addMessage(sessionId, {
    uid: senderUid,
    nickname: `User${senderUid}`, // will be resolved via HTTP API
    headshot: '1',
    content,
    time: now,
    isSelf: senderUid === uid,
  });

  updateSessionLastMessageTime(sessionId, now);
  sortJoinedSessionsByRecentMessage();
  incrementUnread(sessionId);
}

async function fetchSessionTheme(sessionId) {
  try {
    const resp = await fetch(`/sessions/${sessionId}`);
    if (!resp.ok) return;
    const data = await resp.json();
    if (data.theme) {
      useSessionStore.getState().updateSessionTheme(sessionId, data.theme);
    }
  } catch {
    // silent
  }
}

export function sendJoinSession(sessionId) {
  const { socket, connected } = useWsStore.getState();
  if (!socket || !connected) return;

  pendingCommand = 'join';
  socket.send(`join\x1e${sessionId}\x1e`);
}

export function sendCreateSession(theme) {
  const { socket, connected } = useWsStore.getState();
  if (!socket || !connected) return;

  pendingCommand = 'create';
  pendingCreateTheme = theme;
  socket.send(`join\x1e0\x1e${theme}`);
}

export function sendExitSession(sessionId) {
  const { socket, connected } = useWsStore.getState();
  if (!socket || !connected) return;

  pendingCommand = 'exit';
  pendingExitId = sessionId;
  socket.send(`exit\x1e${sessionId}\x1e`);
}

export function sendMessage(sessionId, content) {
  const { socket, connected } = useWsStore.getState();
  if (!socket || !connected) return;

  pendingCommand = 'send';
  pendingSend = { sessionId, content };
  socket.send(`send\x1e${sessionId}\x1e${content}`);
}

export function disconnect() {
  const { socket } = useWsStore.getState();
  if (socket) {
    socket.close();
  }
}
