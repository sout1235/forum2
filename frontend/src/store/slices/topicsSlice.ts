import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Topic } from '../../types';

interface TopicsState {
  topics: Topic[];
  currentTopic: Topic | null;
  loading: boolean;
  error: string | null;
}

const initialState: TopicsState = {
  topics: [],
  currentTopic: null,
  loading: false,
  error: null,
};

const topicsSlice = createSlice({
  name: 'topics',
  initialState,
  reducers: {
    setTopics: (state, action: PayloadAction<Topic[]>) => {
      state.topics = action.payload;
      state.loading = false;
      state.error = null;
    },
    setCurrentTopic: (state, action: PayloadAction<Topic>) => {
      state.currentTopic = action.payload;
      state.loading = false;
      state.error = null;
    },
    addTopic: (state, action: PayloadAction<Topic>) => {
      state.topics.unshift(action.payload);
    },
    updateTopic: (state, action: PayloadAction<Topic>) => {
      const index = state.topics.findIndex(topic => topic.id === action.payload.id);
      if (index !== -1) {
        state.topics[index] = action.payload;
      }
      if (state.currentTopic?.id === action.payload.id) {
        state.currentTopic = action.payload;
      }
    },
    deleteTopic: (state, action: PayloadAction<number>) => {
      state.topics = state.topics.filter(topic => topic.id !== action.payload);
      if (state.currentTopic?.id === action.payload) {
        state.currentTopic = null;
      }
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string>) => {
      state.error = action.payload;
      state.loading = false;
    },
  },
});

export const {
  setTopics,
  setCurrentTopic,
  addTopic,
  updateTopic,
  deleteTopic,
  setLoading,
  setError,
} = topicsSlice.actions;

export default topicsSlice.reducer; 