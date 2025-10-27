'use client'
import {useState, useRef } from "react";
import {Button} from "@/components/ui/button";
import { PlusIcon } from "lucide-react";
import {addFolder } from '@/app/api/folder';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription
} from "@/components/ui/dialog";
import { uploadFiles } from "@/app/api/files";
import { createProject, createTodo, Project } from "../api/todo";


export default function TodoDialog({projects}:{projects:Project[]}) {
    const [newProject, setNewProject] = useState<string>(); 
    const [modalType, setModalType] = useState<"project" | "todo" | null>(null);
    const [dialogOpen, setDialogOpen] = useState(false);
    const [newTodo, setNewTodo] = useState<string>("");
    const [selectedProject, setSelectedProject] = useState<number>();
    const [uploading, setUploading] = useState<boolean>(false);

    const handleCreateProject = async() => {
        try {
          if (newProject!=='' && newProject!==undefined) {
            await createProject(newProject, '');
          }
            alert("Project created successfully");
        }  catch (err) {
            console.error("Error creating project:", err);
            alert("Error creating project");
        }
    }

    const handleCreateTodo = async() => {
        try {
          if (newTodo!=='' && newTodo!==undefined) {
            await createTodo(newTodo, '', undefined);
          }
            alert("Todo created successfully");
        }  catch (err) {
            console.error("Error creating todo:", err);
            alert("Error creating todo");
        }
    }

    
    return (
        <>
        <div className="fixed bottom-6 right-6">
        <Button
          onClick={() => setDialogOpen(true)}
          className="rounded-full h-14 w-14 p-0 shadow-lg"
        >
          <PlusIcon className="h-6 w-6" />
        </Button>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                        <DialogContent className="sm:max-w-md">
                          <DialogHeader>
                            <DialogTitle>Select an option</DialogTitle>
                            <DialogDescription>
                              Choose what you want to upload.
                            </DialogDescription>
                          </DialogHeader>
                          <div className="flex flex-col gap-3 mt-4">
                            <Button
                              onClick={() => {
                                setDialogOpen(false);
                                setModalType("project");
                              }}
                              variant="outline"
                            >
                              Create New Project
                            </Button>
                            <Button
                              onClick={() => {
                                setDialogOpen(false);
                                setModalType("todo");
                              }}
                              variant="outline"
                            >
                              Create Todo
                            </Button>
                          </div>
                        </DialogContent>
                  </Dialog>
        
                  <Dialog open={modalType === "project"} onOpenChange={() => setModalType(null)}>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Create New Project</DialogTitle>
                    </DialogHeader>
                    <form className="flex flex-col gap-4">
                      <input type="text" multiple value={newProject} onChange={(e) => setNewProject(e.target.value)}  className="border rounded p-2" />

                      <Button type="submit" onClick={() => handleCreateProject()}>Submit</Button>
                    </form>
                  </DialogContent>
                </Dialog>
        
                  <Dialog open={modalType === "todo"} onOpenChange={() => setModalType(null)}>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Create Todo</DialogTitle>
                    </DialogHeader>
                    <form className="flex flex-col gap-4">
                      <input type="text" multiple value={newTodo} onChange={(e) => setNewTodo(e.target.value)}  className="border rounded p-2" />
                      <select value={selectedProject} onChange={(e) => setSelectedProject(Number(e.target.value))}>
                        <option value="">Select Project</option>
                        {projects.map((p) => (
                          <option key={p.id} value={p.id}>{p.title}</option>
                        ))}
                      </select>
                      <Button type="submit" onClick={() => handleCreateTodo()}>Submit</Button>
                    </form>
                  </DialogContent>
        </Dialog>
        </>
    )
}
