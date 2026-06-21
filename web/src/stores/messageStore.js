import { create } from 'zustand';

export const useMessageStore = create((set) => ({
  messages: {},

  ensureSession: (sessionId) =>
    set((state) => {
      if (state.messages[sessionId]) return state;
      return { messages: { ...state.messages, [sessionId]: [] } };
    }),

  addMessage: (sessionId, message) =>
    set((state) => ({
      messages: {
        ...state.messages,
        [sessionId]: [...(state.messages[sessionId] || []), message],
      },
    })),

  clearAll: () => set({ messages: {} }),
}));
