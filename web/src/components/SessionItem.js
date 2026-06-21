export default function SessionItem({ session, isActive, unread, onSelect, onExit }) {
  return (
    <li
      className={`session-item ${isActive ? 'session-item--active' : ''}`}
      onClick={onSelect}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => e.key === 'Enter' && onSelect()}
    >
      <div className="session-item__info">
        <span className="session-item__theme">
          {session.theme || `会话 ${session.id}`}
        </span>
      </div>
      {unread > 0 && <span className="badge">{unread}</span>}
      <button
        className="btn btn--exit"
        onClick={(e) => {
          e.stopPropagation();
          onExit();
        }}
        title="退出会话"
      >
        ×
      </button>
    </li>
  );
}
