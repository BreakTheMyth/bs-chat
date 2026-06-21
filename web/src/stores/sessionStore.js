import { create } from 'zustand';

export const useSessionStore = create((set, get) => ({
  joinedSessions: [],
  currentSession: null,
  recommendSessions: [],
  unreadCount: {},

  setCurrentSession: (id) =>
    set((state) => {
      const unreadCount = { ...state.unreadCount };
      if (id != null) delete unreadCount[id];
      return { currentSession: id, unreadCount };
    }),

  addJoinedSession: (session) =>
    set((state) => {
      if (state.joinedSessions.find((s) => s.id === session.id)) return state;
      return { joinedSessions: [...state.joinedSessions, session] };
    }),

  removeJoinedSession: (id) =>
    set((state) => ({
      joinedSessions: state.joinedSessions.filter((s) => s.id !== id),
      currentSession: state.currentSession === id ? null : state.currentSession,
    })),

  setRecommendSessions: (sessions) => set({ recommendSessions: sessions }),

  sortJoinedSessionsByRecentMessage: () =>
    set((state) => {
      const sorted = [...state.joinedSessions].sort(
        (a, b) => (b.lastMessageTime || 0) - (a.lastMessageTime || 0)
      );
      return { joinedSessions: sorted };
    }),

  updateSessionLastMessageTime: (id, time) =>
    set((state) => ({
      joinedSessions: state.joinedSessions.map((s) =>
        s.id === id ? { ...s, lastMessageTime: time } : s
      ),
    })),

  incrementUnread: (sessionId) =>
    set((state) => {
      if (state.currentSession === sessionId) return state;
      return {
        unreadCount: {
          ...state.unreadCount,
          [sessionId]: (state.unreadCount[sessionId] || 0) + 1,
        },
      };
    }),

  updateSessionTheme: (id, theme) =>
    set((state) => ({
      joinedSessions: state.joinedSessions.map((s) =>
        s.id === id ? { ...s, theme } : s
      ),
    })),

  clearAll: () =>
    set({
      joinedSessions: [],
      currentSession: null,
      recommendSessions: [],
      unreadCount: {},
    }),
}));
