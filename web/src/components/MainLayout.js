import Sidebar from './Sidebar';
import ChatArea from './ChatArea';
import { useSessionStore } from '../stores/sessionStore';

export default function MainLayout() {
  const currentSession = useSessionStore((s) => s.currentSession);

  return (
    <div className="main-layout">
      <div className="sidebar-panel">
        <Sidebar />
      </div>
      <div className="chat-panel">
        <ChatArea key={currentSession ?? '_empty'} />
      </div>
    </div>
  );
}
