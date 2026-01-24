import { useCallback, useEffect, useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import type { Post, Comment } from "../../types";
import { api, getErrorMessage } from "../../services/api";
import { useAuth } from "../../context/AuthContext";
import { Button } from "../ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Alert, AlertDescription } from "../ui/alert";
import { formatDate } from "../../lib/utils";
import { CommentList } from "../Comments/CommentList";
import { CommentBox } from "../Comments/CommentBox";
import { PostForm } from "./PostForm";

// Breadcrumb navigation component
function Breadcrumb({ post }: { post: Post }) {
  return (
    <nav className="flex items-center space-x-2 text-sm text-gray-600">
      <Link to="/topics" className="hover:text-blue-600">Topics</Link>
      <span>/</span>
      <Link to={`/topics/${post.topic_id}`} className="hover:text-blue-600">
        {post.topic_name}
      </Link>
      <span>/</span>
      <span>{post.title}</span>
    </nav>
  );
}

export function PostDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();

  const [post, setPost] = useState<Post | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [isEditing, setIsEditing] = useState(false);

  const postId = id ? parseInt(id, 10) : NaN;

  const fetchData = useCallback(async () => {
    if (Number.isNaN(postId)) return;

    setIsLoading(true);
    setError("");

    try {
      const [postData, commentsData] = await Promise.all([
        api.getPost(postId),
        api.getComments(postId),
      ]);
      setPost(postData);
      setComments(commentsData ?? []);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  }, [postId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleDelete = async () => {
    if (!post || !confirm("Are you sure you want to delete this post?")) return;

    try {
      await api.deletePost(post.id);
      navigate(`/topics/${post.topic_id}`);
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleCommentAdded = (comment: Comment) => {
    setComments((prev) => [...prev, comment]);
  };

  const handleCommentUpdated = (updated: Comment) => {
    setComments((prev) => prev.map((c) => (c.id === updated.id ? updated : c)));
  };

  const handleCommentDeleted = (commentId: number) => {
    setComments((prev) => prev.filter((c) => c.id !== commentId));
  };

  const handleEditSuccess = () => {
    setIsEditing(false);
    fetchData();
  };

  // Loading state
  if (isLoading) {
    return <div className="text-center py-8">Loading...</div>;
  }

  // Not found state
  if (!post) {
    return <div className="text-center py-8">Post not found</div>;
  }

  // Edit mode
  if (isEditing) {
    return (
      <PostForm
        topicId={post.topic_id}
        post={post}
        onSuccess={handleEditSuccess}
        onCancel={() => setIsEditing(false)}
      />
    );
  }

  const isOwner = user?.id === post.user_id;
  const wasEdited = post.updated_at !== post.created_at;

  return (
    <div className="space-y-6">
      <Breadcrumb post={post} />

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div>
              <CardTitle>{post.title}</CardTitle>
              <CardDescription className="mt-1">
                by {post.username} • {formatDate(post.created_at)}
                {wasEdited && " (edited)"}
              </CardDescription>
            </div>
            {isOwner && (
              <div className="flex space-x-2">
                <Button variant="outline" size="sm" onClick={() => setIsEditing(true)}>
                  Edit
                </Button>
                <Button variant="destructive" size="sm" onClick={handleDelete}>
                  Delete
                </Button>
              </div>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-gray-700 whitespace-pre-wrap">{post.content}</p>
        </CardContent>
      </Card>

      <section className="space-y-4">
        <h2 className="text-xl font-semibold">Comments ({comments.length})</h2>
        <CommentBox postId={post.id} onCommentAdded={handleCommentAdded} />
        <CommentList
          comments={comments}
          onCommentUpdated={handleCommentUpdated}
          onCommentDeleted={handleCommentDeleted}
        />
      </section>
    </div>
  );
}
