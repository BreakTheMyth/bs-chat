import { create } from 'zustand';

export const useWsStore = create((set) => ({
  socket: null,
  connected: false,
  connecting: false,

  setSocket: (socket) => set({ socket }),
  setConnected: (connected) => set({ connected }),
  setConnecting: (connecting) => set({ connecting }),

  clearSocket: () => set({ socket: null, connected: false, connecting: false }),
}));
