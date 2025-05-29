import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface Comment {
  id: number;
  content: string;
  authorId: number;
  author: {
    id: number;
    username: string;
    avatar?: string;
  };
  topicId: number;
  parentId?: number;
  replies: Comment[];
  likes: number;
  createdAt: string;
  updatedAt: string;
}

interface CommentsState {
  comments: Comment[];
  loading: boolean;
  error: string | null;
}

const initialState: CommentsState = {
  comments: [],
  loading: false,
  error: null,
};

const commentsSlice = createSlice({
  name: 'comments',
  initialState,
  reducers: {
    fetchCommentsStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    fetchCommentsSuccess: (state, action: PayloadAction<Comment[]>) => {
      state.loading = false;
      state.comments = action.payload;
    },
    fetchCommentsFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    addComment: (state, action: PayloadAction<Comment>) => {
      if (action.payload.parentId) {
        const parentComment = state.comments.find(c => c.id === action.payload.parentId);
        if (parentComment) {
          parentComment.replies.push(action.payload);
        }
      } else {
        state.comments.push(action.payload);
      }
    },
    updateComment: (state, action: PayloadAction<Comment>) => {
      const updateCommentInArray = (comments: Comment[]): boolean => {
        for (let i = 0; i < comments.length; i++) {
          if (comments[i].id === action.payload.id) {
            comments[i] = action.payload;
            return true;
          }
          if (updateCommentInArray(comments[i].replies)) {
            return true;
          }
        }
        return false;
      };
      updateCommentInArray(state.comments);
    },
    deleteComment: (state, action: PayloadAction<number>) => {
      const deleteCommentFromArray = (comments: Comment[]): boolean => {
        for (let i = 0; i < comments.length; i++) {
          if (comments[i].id === action.payload) {
            comments.splice(i, 1);
            return true;
          }
          if (deleteCommentFromArray(comments[i].replies)) {
            return true;
          }
        }
        return false;
      };
      deleteCommentFromArray(state.comments);
    },
  },
});

export const {
  fetchCommentsStart,
  fetchCommentsSuccess,
  fetchCommentsFailure,
  addComment,
  updateComment,
  deleteComment,
} = commentsSlice.actions;

export default commentsSlice.reducer; 