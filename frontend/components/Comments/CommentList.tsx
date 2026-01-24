import { useState } from "react";
import type { Comment } from "../../types";
import { api, getErrorMessage } from "../../services/api";
import { useAuth } from "../../context/AuthContext";
import { Button } from "../../components/ui/button";
import { Textarea } from "../../components/ui/textarea";
import { Card, CardContent } from "../../components/ui/card";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { formatDate } from "../../lib/utils";

interface CommentListProps {
  comments: Comment[];
  onCommentUpdated: (comment: Comment) => void;
  onCommentDeleted: (commentId: number) => void;
}

export function CommentList({
  comments,
  onCommentUpdated,
  onCommentDeleted,
}: CommentListProps) {
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editContent, setEditContent] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { user } = useAuth();

  const handleEdit = (comment: Comment) => {
    setEditingId(comment.id);
    setEditContent(comment.content);
    setError("");
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setEditContent("");
    setError("");
  };

  const handleSaveEdit = async (commentId: number) => {
    if (!editContent.trim()) {
      setError("Comment cannot be empty");
      return;
    }

    setIsLoading(true);
    try {
      const updatedComment = await api.updateComment(commentId, {
        content: editContent,
      });
      onCommentUpdated(updatedComment);
      setEditingId(null);
      setEditContent("");
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (commentId: number) => {
    if (!confirm("Are you sure you want to delete this comment?")) return;

    try {
      await api.deleteComment(commentId);
      onCommentDeleted(commentId);
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  if (comments.length === 0) {
    return (
      <p className="text-gray-500 text-center py-4">
        No comments yet. Be the first to comment!
      </p>
    );
  }

  return (
    <div className="space-y-4">
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {comments.map((comment) => (
        <Card key={comment.id}>
          <CardContent className="pt-4">
            {editingId === comment.id ? (
              <div className="space-y-3">
                <Textarea
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  disabled={isLoading}
                />
                <div className="flex space-x-2">
                  <Button
                    size="sm"
                    onClick={() => handleSaveEdit(comment.id)}
                    disabled={isLoading}
                  >
                    {isLoading ? "Saving..." : "Save"}
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={handleCancelEdit}
                    disabled={isLoading}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            ) : (
              <>
                <div className="flex items-start justify-between">
                  <div className="text-sm text-gray-600 mb-2">
                    <span className="font-medium">{comment.username}</span>
                    <span className="mx-1">•</span>
                    <span>{formatDate(comment.created_at)}</span>
                    {comment.updated_at !== comment.created_at && (
                      <span className="text-gray-400"> (edited)</span>
                    )}
                  </div>
                  {user?.id === comment.user_id && (
                    <div className="flex space-x-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleEdit(comment)}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="text-red-600 hover:text-red-700"
                        onClick={() => handleDelete(comment.id)}
                      >
                        Delete
                      </Button>
                    </div>
                  )}
                </div>
                <p className="text-gray-700 whitespace-pre-wrap">{comment.content}</p>
              </>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}