import { useState } from 'react';

export default function CreateSessionModal({ onConfirm, onCancel }) {
  const [theme, setTheme] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    const trimmed = theme.trim();
    if (!trimmed) return;
    onConfirm(trimmed);
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h3>创建会话</h3>
        <form onSubmit={handleSubmit}>
          <input
            className="input"
            type="text"
            value={theme}
            onChange={(e) => setTheme(e.target.value)}
            placeholder="输入会话名..."
            maxLength={200}
            autoFocus
          />
          <div className="btn-row">
            <button className="btn btn--primary" type="submit">
              创建
            </button>
            <button className="btn" type="button" onClick={onCancel}>
              取消
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
