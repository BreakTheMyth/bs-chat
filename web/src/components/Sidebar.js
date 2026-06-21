import { useState, useCallback, useEffect } from 'react';
import { useSessionStore } from '../stores/sessionStore';
import { fetchRecommendedSessions } from '../services/api';
import {
  sendJoinSession,
  sendCreateSession,
  sendExitSession,
} from '../services/websocket';
import SessionItem from './SessionItem';
import CreateSessionModal from './CreateSessionModal';

export default function Sidebar() {
  const joinedSessions = useSessionStore((s) => s.joinedSessions);
  const currentSession = useSessionStore((s) => s.currentSession);
  const recommendSessions = useSessionStore((s) => s.recommendSessions);
  const setRecommendSessions = useSessionStore((s) => s.setRecommendSessions);
  const unreadCount = useSessionStore((s) => s.unreadCount);

  const [joinId, setJoinId] = useState('');
  const [showCreate, setShowCreate] = useState(false);

  const loadRecommendations = useCallback(async () => {
    if (joinedSessions.length > 0) return;
    try {
      const data = await fetchRecommendedSessions();
      setRecommendSessions(Array.isArray(data) ? data : []);
    } catch {
      // silent
    }
  }, [joinedSessions.length, setRecommendSessions]);

  useEffect(() => {
    loadRecommendations();
  }, [loadRecommendations]);

  const handleJoin = () => {
    const id = parseInt(joinId, 10);
    if (isNaN(id) || id <= 0) return;
    sendJoinSession(id);
    setJoinId('');
  };

  const handleJoinRecommended = (id) => {
    sendJoinSession(id);
  };

  const handleCreate = (theme) => {
    sendCreateSession(theme);
    setShowCreate(false);
  };

  const handleExit = (id) => {
    sendExitSession(id);
  };

  return (
    <div className="sidebar">
      <div className="sidebar-section">
        <div className="sidebar-header">
          <h3 className="sidebar-title">已加入会话</h3>
        </div>

        {joinedSessions.length > 0 ? (
          <ul className="session-list">
            {joinedSessions.map((s) => (
              <SessionItem
                key={s.id}
                session={s}
                isActive={currentSession === s.id}
                unread={unreadCount[s.id] || 0}
                onSelect={() =>
                  useSessionStore.getState().setCurrentSession(s.id)
                }
                onExit={() => handleExit(s.id)}
              />
            ))}
          </ul>
        ) : null}

        <button className="btn btn--create" onClick={() => setShowCreate(true)}>
          + 创建会话
        </button>
      </div>

      {joinedSessions.length === 0 && recommendSessions.length > 0 && (
        <div className="sidebar-section">
          <div className="sidebar-header">
            <h3 className="sidebar-title">推荐会话</h3>
            <button className="btn btn--small" onClick={loadRecommendations}>
              刷新
            </button>
          </div>
          <ul className="session-list">
            {recommendSessions.map((s) => (
              <li key={s.id} className="session-item">
                <div className="session-item__info">
                  <span className="session-item__theme">{s.theme}</span>
                  <span className="session-item__online">{s.online} 在线</span>
                </div>
                <button
                  className="btn btn--small"
                  onClick={() => handleJoinRecommended(s.id)}
                >
                  加入
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="sidebar-section">
        <div className="sidebar-header">
          <h3 className="sidebar-title">加入会话</h3>
        </div>
        <div className="join-row">
          <input
            className="input input--small"
            type="text"
            value={joinId}
            onChange={(e) => setJoinId(e.target.value)}
            placeholder="输入 Session ID"
            onKeyDown={(e) => e.key === 'Enter' && handleJoin()}
          />
          <button className="btn btn--small" onClick={handleJoin}>
            加入
          </button>
        </div>
      </div>

      {showCreate && (
        <CreateSessionModal
          onConfirm={handleCreate}
          onCancel={() => setShowCreate(false)}
        />
      )}
    </div>
  );
}
