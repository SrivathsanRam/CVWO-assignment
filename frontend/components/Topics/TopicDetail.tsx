import React, { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import type { Topic, Post } from "../../types";
import { api, getErrorMessage } from "../../services/api";
import { useAuth } from "../../context/AuthContext";
import { Button } from "../../components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { formatDate } from "../../lib/utils";
import { PostForm } from "../Posts/PostForm";

export function TopicDetail() {
  const { id } = useParams<{ id: string }>();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [showPostForm, setShowPostForm] = useState(false);
  const { user } = useAuth();

  const fetchData = async () => {
    if (!id) return;

    try {
      const [topicData, postsData] = await Promise.all([
        api.getTopic(parseInt(id)),
        api.getPosts(parseInt(id)),
      ]);
      setTopic(topicData);
      setPosts(postsData || []);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [id]);

  const handleDeletePost = async (postId: number) => {
    if (!confirm("Are you sure you want to delete this post?")) return;

    try {
      await api.deletePost(postId);
      setPosts(posts.filter((p) => p.id !== postId));
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handlePostSuccess = () => {
    setShowPostForm(false);
    fetchData();
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading...</div>;
  }

  if (!topic) {
    return <div className="text-center py-8">Topic not found</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center space-x-2 text-sm text-gray-600">
        <Link to="/topics" className="hover:text-blue-600">
          Topics
        </Link>
        <span>/</span>
        <span>{topic.title}</span>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{topic.title}</CardTitle>
          <CardDescription>
            Created by {topic.username} • {formatDate(topic.created_at)}
          </CardDescription>
        </CardHeader>
        {topic.description && (
          <CardContent>
            <p className="text-gray-600">{topic.description}</p>
          </CardContent>
        )}
      </Card>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Posts</h2>
        <Button onClick={() => setShowPostForm(true)}>Create Post</Button>
      </div>

      {showPostForm && (
        <PostForm
          topicId={topic.id}
          onSuccess={handlePostSuccess}
          onCancel={() => setShowPostForm(false)}
        />
      )}

      {posts.length === 0 ? (
        <Card>
          <CardContent className="py-8 text-center text-gray-500">
            No posts in this topic yet. Be the first to create one!
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
                      by {post.username} • {formatDate(post.created_at)}
                    </CardDescription>
                  </div>
                  {user?.id === post.user_id && (
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDeletePost(post.id)}
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