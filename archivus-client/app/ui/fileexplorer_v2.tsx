'use client'
import { useEffect,useState, useRef } from "react";
import { FileMetaData, getFilesByFolder, getFilesByFolderResponse } from "../api/files";
import FileCard from "./components/filecard";
import { DndContext, closestCenter,DragOverlay, useSensors, useSensor, MouseSensor, TouchSensor, DragEndEvent,DragStartEvent, PointerSensor } from '@dnd-kit/core';
import { useSortable } from '@dnd-kit/sortable';
import { moveFile } from "@/app/api/files";
import { uploadFiles } from "@/app/api/files";
import { CSS } from '@dnd-kit/utilities';
import { useSearchParams } from "next/navigation";
import {
    SortableContext,
    verticalListSortingStrategy,
} from '@dnd-kit/sortable';

import FileFolderDialog from "./components/fileFolderModal";
import PathBreadcrumbs from "./components/breadcrumbs";

const FileDraggableCard: React.FC<{
    file: FileMetaData;
    onClick: (file: FileMetaData) => void;
}> = ({file, onClick}) => {
    const {attributes, listeners, setNodeRef, transform, transition, isDragging} = useSortable({id: file.id});
    const style:React.CSSProperties  = {
        transform: CSS.Transform.toString(transform),
        transition: transition || "transform 200ms ease", // 👈 fallback
        opacity: isDragging ? 0.6 : 1,
        cursor: isDragging ? "grabbing" : "grab",        
        // pointerEvents: 'auto'
    };
    const [dragging,setDragging] = useState(false);
    const handleMouseDown = () => setDragging(false);
    const handleMouseMove = () => setDragging(true);
    const handleMouseUp = () => {
        if (!dragging) onClick(file);
    };
    return (
        <div 
        ref={setNodeRef}
        style={style}
        {...attributes}
        {...listeners}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        >
            <FileCard file={file} />
        </div>
    )

}
const generateRandomId = () => {
    return "gid_"+Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
}



export default function FileExplorer() {
      const searchParams = useSearchParams();
      const folder = searchParams.get('folder') || '';
        const [files, setFiles] = useState<Map<string,FileMetaData>>(new Map<string, FileMetaData>());
        const [loading, setLoading] = useState<boolean>(false);
        const [size, setSize] = useState<number>(0);
        const [activeId, setActiveId] = useState<string | null>(null);
        const parts = folder.split("/").filter(Boolean);


        useEffect(() => {
        async function loadFiles() {
                  try {
                    setLoading(true);
                    const data = await getFilesByFolder(folder);
                    for (const file of data.files) {
                      if (file.IsDir || !file.ID || file.ID === '') {
                        file.ID = generateRandomId();
                      }
                      file.id = generateRandomId();
                    }
                    const fileMap = new Map<string,FileMetaData>();
                    data.files.forEach((f) => fileMap!.set(f.id,f));
                    setFiles(fileMap);
                    setSize(data.size);
                    console.log(fileMap);
                  } catch (err) {
                    console.error("Error fetching files", err);
                  } finally {
                    setLoading(false);
                  }
                }
        loadFiles();
      },[folder]);
      
      const sensors = useSensors(
        useSensor(PointerSensor, {activationConstraint: {distance: 5}}),
        useSensor(MouseSensor));

    const handleDragStart = (event: DragStartEvent) => {
        // console.log("drag start", event.active.id);
        setActiveId(event.active.id.toString());
    };
    const handleDragEnd = (event: DragEndEvent) => {
        const {active, over} = event;
        // console.log("Drag ended:", {active, over});
        if (!over ||active.id == over.id) return;
        const activeFile = files.get(active.id.toString());
        if (activeFile && activeFile.IsDir) {
          alert("Cannot move folders");
          setActiveId(null);
          return;
        }

        const overFile = files.get(over.id.toString());
        if (overFile?.IsDir) {
          moveFile(activeFile!.Path, overFile.Path).then((res) => {
            if (res) {
              alert("File moved successfully");
              window.location.reload();
            }
          }).catch((err) => {
            console.error("Error moving file", err);
            alert("Error moving file");
          });
          setActiveId(null);
          return;
        }

    }
    const handleClick = (file: FileMetaData) => {
        if (file.IsDir) {
            window.location.href = `?folder=${file.Path}`;
        } else {
            window.open(file.SignedUrl, '_blank');
        }
    }

      return (
        <>
        <div className="flex mx-auto p-4  justify-between ">
          <PathBreadcrumbs parts={parts} />
        </div>
        <FileFolderDialog folder={folder} />
        <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <SortableContext items={Array.from(files!.values())} strategy={verticalListSortingStrategy}>
        <ul role="list" className="grid grid-cols-6 gap-6 sm:grid-cols-4 lg:grid-cols-7">
          {Array.from(files!.values()).map((file) => (
            <FileDraggableCard key={file.id} file={file} onClick={handleClick} />
          ))}
        </ul>
          <DragOverlay>
                {activeId ? (
                  <div className="p-3 bg-white border rounded-md shadow-lg scale-105 opacity-90">
                    <FileDraggableCard key={activeId} file={files!.get(activeId)!} onClick={() => {}}/>
                  </div>
                ) : null}
          </DragOverlay>
        
      </SortableContext>
    </DndContext>
    </>
      )
}