import { useState, useRef, useCallback, useEffect } from 'react';
import { sendMessage } from '../services/websocket';

const MAX_LENGTH = 0x990; // 2448

export default function MessageInput({ sessionId }) {
  const [text, setText] = useState('');
  const textareaRef = useRef(null);

  // Auto-resize textarea
  useEffect(() => {
    const el = textareaRef.current;
    if (!el) return;
    el.style.height = 'auto';
    el.style.height = Math.min(el.scrollHeight, 200) + 'px';
  }, [text]);

  const handleSend = useCallback(() => {
    const trimmed = text.trim();
    if (!trimmed || trimmed.length > MAX_LENGTH) return;

    sendMessage(sessionId, trimmed);

    setText('');
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
    }
  }, [text, sessionId]);

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const remaining = MAX_LENGTH - text.length;

  return (
    <div className="input-area">
      <textarea
        ref={textareaRef}
        className="input-area__textarea"
        value={text}
        onChange={(e) => setText(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="输入消息... (Enter 发送, Shift+Enter 换行)"
        maxLength={MAX_LENGTH}
        rows={1}
      />
      <div className="input-area__footer">
        <span className={`input-area__counter ${remaining < 100 ? 'input-area__counter--warn' : ''}`}>
          {remaining}
        </span>
        <button
          className="btn btn--primary"
          onClick={handleSend}
          disabled={!text.trim()}
        >
          发送
        </button>
      </div>
    </div>
  );
}
