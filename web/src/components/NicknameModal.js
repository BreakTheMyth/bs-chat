import { useState, useCallback, useEffect, useRef } from 'react';
import { useUserStore } from '../stores/userStore';
import { useWsStore } from '../stores/wsStore';
import { connectWebSocket } from '../services/websocket';

const HEADSHOTS = ['1', '2', '3', '4', '5'];
const TIMEOUT_MS = 10000;

function hasControlChar(text) {
  for (let i = 0; i < text.length; i++) {
    const code = text.charCodeAt(i);
    if (code <= 31 || code === 127) return true;
  }
  return false;
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.appendChild(document.createTextNode(text));
  return div.innerHTML;
}

export default function NicknameModal() {
  const { nickname: savedNickname, headshot: savedHeadshot, setLocalInfo } =
    useUserStore();
  const wsConnecting = useWsStore((s) => s.connecting);
  const wsConnected = useWsStore((s) => s.connected);

  const [nickname, setNickname] = useState(savedNickname);
  const [headshot, setHeadshot] = useState(savedHeadshot);
  const [error, setError] = useState('');
  const [submitted, setSubmitted] = useState(false);
  const submittedRef = useRef(false);

  // Sync ref so the timeout callback can read latest submitted state
  useEffect(() => {
    submittedRef.current = submitted;
  }, [submitted]);

  // When WebSocket closes unexpectedly after submit, show an error
  useEffect(() => {
    if (submitted && !wsConnecting && !wsConnected) {
      setSubmitted(false);
      setError('连接失败，请重试');
    }
  }, [submitted, wsConnecting, wsConnected]);

  const handleSubmit = useCallback(
    (e) => {
      e.preventDefault();
      const trimmed = nickname.trim();

      if (!trimmed) {
        setError('昵称不能为空');
        return;
      }
      if (trimmed.length > 32) {
        setError('昵称最多32个字符');
        return;
      }
      if (hasControlChar(trimmed)) {
        setError('昵称不能包含控制字符');
        return;
      }

      setError('');
      setLocalInfo(trimmed, headshot);
      setSubmitted(true);

      const socket = connectWebSocket();
      if (!socket) {
        setSubmitted(false);
        setError('连接失败，请重试');
        return;
      }

      // Safety timeout: if we don't get a UID within TIMEOUT_MS, show error
      const tid = setTimeout(() => {
        if (submittedRef.current) {
          setSubmitted(false);
          setError('连接超时，请重试');
        }
      }, TIMEOUT_MS);

      socket.addEventListener('open', () => clearTimeout(tid));
      socket.addEventListener('close', () => clearTimeout(tid));
    },
    [nickname, headshot, setLocalInfo]
  );

  const isConnecting = submitted || wsConnecting;

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h1 className="modal-title">BSChat</h1>
        <p className="modal-subtitle">无需注册，即开即用</p>

        <form onSubmit={handleSubmit}>
          <label className="field-label">昵称</label>
          <input
            className="input"
            type="text"
            value={nickname}
            onChange={(e) => {
              setNickname(e.target.value);
              setError('');
            }}
            placeholder="输入昵称..."
            maxLength={32}
            autoFocus
            disabled={isConnecting}
          />

          <label className="field-label">选择头像</label>
          <div className="headshot-grid">
            {HEADSHOTS.map((h) => (
              <img
                key={h}
                src={`/headshot/${h}`}
                alt={`头像${h}`}
                className={`headshot ${headshot === h ? 'headshot--selected' : ''}`}
                onClick={() => !isConnecting && setHeadshot(h)}
              />
            ))}
          </div>

          {error && <div className="error-msg">{escapeHtml(error)}</div>}

          <button
            className="btn btn--primary btn--full"
            type="submit"
            disabled={isConnecting}
          >
            {isConnecting ? '连接中...' : '进入聊天室'}
          </button>
        </form>
      </div>
    </div>
  );
}
