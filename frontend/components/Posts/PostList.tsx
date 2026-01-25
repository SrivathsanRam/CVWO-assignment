import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import type { Post } from "../../types";
import { api, getErrorMessage } from "../../services/api";
import { useAuth } from "../../context/AuthContext";
import { Button } from "../../components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { formatDate } from "../../lib/utils";

export function PostList() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const { user } = useAuth();

  const fetchPosts = async () => {
    try {
      const data = await api.getPosts();
      setPosts(data || []);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this post?")) return;

    try {
      await api.deletePost(id);
      setPosts(posts.filter((p) => p.id !== id));
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading posts...</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">All Posts</h1>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {posts.length === 0 ? (
        <Card>
          <CardContent className="py-8 text-center text-gray-500">
            No posts yet. Go to a topic to create one!
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {posts.map((post) => (
            <Card key={post.id} className="hover:shadow-md transition-shadow">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <Link to={`/posts/${post.id}`}>
                      <CardTitle className="text-lg hover:text-blue-600 cursor-pointer">
                        {post.title}
                      </CardTitle>
                    </Link>
                    <CardDescription className="mt-1">
                      in{" "}
                      <Link
                        to={`/topics/${post.topic_id}`}
                        className="text-blue-600 hover:underline"
                      >
                        {post.topic_name}
                      </Link>{" "}
                      • by {post.username} • {formatDate(post.created_at)}
                    </CardDescription>
                  </div>
                  {user?.id === post.user_id && (
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDelete(post.id)}
                    >
                      Delete
                    </Button>
                  )}
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600 line-clamp-3">{post.content}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}