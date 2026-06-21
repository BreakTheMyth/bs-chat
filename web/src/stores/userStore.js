import { create } from 'zustand';

const NICKNAME_KEY = 'bschat_nickname';
const HEADSHOT_KEY = 'bschat_headshot';

export const useUserStore = create((set) => ({
  uid: null,
  nickname: localStorage.getItem(NICKNAME_KEY) || '',
  headshot: localStorage.getItem(HEADSHOT_KEY) || '1',

  setUserInfo: (uid, nickname, headshot) => {
    localStorage.setItem(NICKNAME_KEY, nickname);
    localStorage.setItem(HEADSHOT_KEY, headshot);
    set({ uid, nickname, headshot });
  },

  setLocalInfo: (nickname, headshot) => {
    localStorage.setItem(NICKNAME_KEY, nickname);
    localStorage.setItem(HEADSHOT_KEY, headshot);
    set({ nickname, headshot });
  },

  clearUser: () => {
    localStorage.removeItem(NICKNAME_KEY);
    localStorage.removeItem(HEADSHOT_KEY);
    set({ uid: null, nickname: '', headshot: '1' });
  },
}));
