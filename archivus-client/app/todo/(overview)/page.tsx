'use client';

import { useEffect, useState } from 'react';
import { getProjects,createProject, Project, Todo , createTodo, getTodos} from '../../api/todo';
import TodoTable from '@/app/ui/todoTable';
import ProjectTable from '@/app/ui/projectTable';
import TodoDialog from '@/app/ui/notesModal';



export default function Page() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [todos, setTodos] = useState<Todo[]>([]);
  const [projectName, setProjectName] = useState('');
  const [todoTitle, setTodoTitle] = useState('');
  const [selectedProject, setSelectedProject] = useState<number>();

  useEffect(() => {
    getProjects().then((res) => {
      setProjects(res);
    });
    getTodos(selectedProject).then((res) => {
      setTodos(res);
    });
  }, [selectedProject]);

  // Handle project creation
  const handleCreateProject = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    createProject(projectName, '').then((res) => {
      setProjects([...projects, res]);
      setProjectName('');
    });
    window.location.reload();
  };

  // Handle todo creation
  const handleCreateTodo = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    createTodo(todoTitle, '', selectedProject ).then((res) => {
      setTodos([...todos, res]);
      setTodoTitle('');
      setSelectedProject(undefined);
    });
  };
  
  const handleProjectChange = (projectId?: number) => {
    if (!projectId) { 
      setSelectedProject(undefined) 
    } else {
      setSelectedProject(projectId);
    }
    console.log(selectedProject);
  }

  return (
<div className="min-h-screen bg-gray-100 p-4 sm:p-6">
      <div className="mx-auto max-w-4xl">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-800 mb-6 text-center">
          Todos
        </h1>
        <TodoDialog projects={projects}/>


        {/* Display Projects */}
        <div className="mb-6">
          <h2 className="text-lg sm:text-xl font-semibold text-gray-700 mb-4">
            Projects
          </h2>
          {projects.length === 0 ? (
            <p className="text-gray-500">No projects yet.</p>
          ) : (
            <ProjectTable projects={projects} onChange={handleProjectChange}/>
          )}
        </div>

        {/* Display Todos */}
        <div>
          <h2 className="text-lg sm:text-xl font-semibold text-gray-700 mb-4">
            Todos
          </h2>
          <TodoTable todos={todos}/>
        </div>
      </div>
    </div>
  );
}