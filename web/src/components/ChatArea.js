import { useRef, useEffect } from 'react';
import { useSessionStore } from '../stores/sessionStore';
import { useMessageStore } from '../stores/messageStore';
import MessageBubble from './MessageBubble';
import MessageInput from './MessageInput';
import EmptyState from './EmptyState';

export default function ChatArea() {
  const currentSessionId = useSessionStore((s) => s.currentSession);
  const joinedSessions = useSessionStore((s) => s.joinedSessions);
  const messages = useMessageStore((s) => s.messages);
  const bottomRef = useRef(null);

  const currentSession = currentSessionId != null
    ? joinedSessions.find((s) => s.id === currentSessionId)
    : null;
  const currentMessages = currentSessionId != null ? messages[currentSessionId] || [] : [];

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [currentMessages.length]);

  if (currentSessionId == null) {
    return <EmptyState />;
  }

  return (
    <div className="chat-area">
      <div className="chat-header">
        <h2>
          {currentSession?.theme || `会话 ${currentSessionId}`}
          <span className="chat-header__id">#{currentSessionId}</span>
        </h2>
      </div>
      <div className="message-list">
        {currentMessages.length === 0 && (
          <p className="message-list__hint">暂无消息，发送第一条消息吧</p>
        )}
        {currentMessages.map((msg, i) => (
          <MessageBubble key={i} message={msg} />
        ))}
        <div ref={bottomRef} />
      </div>
      <MessageInput sessionId={currentSessionId} />
    </div>
  );
}
