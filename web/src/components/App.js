import { useUserStore } from '../stores/userStore';
import NicknameModal from './NicknameModal';
import MainLayout from './MainLayout';

export default function App() {
  const uid = useUserStore((s) => s.uid);

  if (uid == null) {
    return <NicknameModal />;
  }

  return <MainLayout />;
}
