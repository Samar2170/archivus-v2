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

const FileDraggableCard: React.FC<{
    file: FileMetaData;
    onClick: (file: FileMetaData) => void;
}> = ({file, onClick}) => {
    const {attributes, listeners, setNodeRef, transform, transition, isDragging} = useSortable({id: file.ID});
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
    return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
}

export default function FileExplorer() {
      const searchParams = useSearchParams();
      const folder = searchParams.get('folder') || '';
        const [files, setFiles] = useState<FileMetaData[]>([]);
        const [loading, setLoading] = useState<boolean>(false);
        const [size, setSize] = useState<number>(0);
        
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
                    setFiles(data.files || []);
                    setSize(data.size);
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
        console.log("drag start", event.active.id);
    };
        const handleDragEnd = (event: DragEndEvent) => {
            const {active, over} = event;
            console.log("Drag ended:", {active, over});
        }
        const handleClick = (file: FileMetaData) => {
            console.log("Clicked file:", file);
        }

      return (
        // <FileFolderDialog folder={folder} />
        <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <SortableContext items={files.map((f) => f.id)} strategy={verticalListSortingStrategy}>
        <div className="flex flex-col gap-3 p-4 w-full max-w-md mx-auto">
          {files.map((file,indx) => (
            <FileDraggableCard key={indx} file={file} onClick={handleClick} />
          ))}
        </div>
      </SortableContext>
    </DndContext>
      )
}