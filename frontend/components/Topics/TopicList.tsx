import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import type { Topic } from "../../types";
import { api, getErrorMessage } from "../../services/api";
import { useAuth } from "../../context/AuthContext";
import { Button } from "../../components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card";
import { Alert, AlertDescription } from "../../components/ui/alert";
import { formatDate } from "../../lib/utils";
import { TopicForm } from "./TopicForm";

export function TopicList() {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [editingTopic, setEditingTopic] = useState<Topic | null>(null);
  const { user } = useAuth();

  const fetchTopics = async () => {
    try {
      const data = await api.getTopics();
      setTopics(data || []);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchTopics();
  }, []);

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this topic?")) return;

    try {
      await api.deleteTopic(id);
      setTopics(topics.filter((t) => t.id !== id));
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleFormSuccess = () => {
    setShowForm(false);
    setEditingTopic(null);
    fetchTopics();
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading topics...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Topics</h1>
        <Button onClick={() => setShowForm(true)}>Create Topic</Button>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {(showForm || editingTopic) && (
        <TopicForm
          topic={editingTopic}
          onSuccess={handleFormSuccess}
          onCancel={() => {
            setShowForm(false);
            setEditingTopic(null);
          }}
        />
      )}

      {topics.length === 0 ? (
        <Card>
          <CardContent className="py-8 text-center text-gray-500">
            No topics yet. Be the first to create one!
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {topics.map((topic) => (
            <Card key={topic.id} className="hover:shadow-md transition-shadow">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <Link to={`/topics/${topic.id}`}>
                      <CardTitle className="text-lg hover:text-blue-600 cursor-pointer">
                        {topic.title}
                      </CardTitle>
                    </Link>
                    <CardDescription className="mt-1">
                      by {topic.username} • {formatDate(topic.created_at)}
                    </CardDescription>
                  </div>
                  {user?.id === topic.user_id && (
                    <div className="flex space-x-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setEditingTopic(topic)}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDelete(topic.id)}
                      >
                        Delete
                      </Button>
                    </div>
                  )}
                </div>
              </CardHeader>
              {topic.description && (
                <CardContent>
                  <p className="text-gray-600">{topic.description}</p>
                </CardContent>
              )}
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}