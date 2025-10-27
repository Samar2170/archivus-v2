import { EnvelopeIcon, PhoneIcon } from '@heroicons/react/20/solid'
import { Project } from '../api/todo'


export default function ProjectTable({projects, onChange}: {projects: Project[], onChange: (id?: number) => void}) { 
  return (
    <ul role="list" className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {projects && projects.map((p) => (
        <li key={p.id} className="col-span-1 divide-y divide-gray-200 rounded-lg bg-white shadow-sm">
        <button onClick={() => onChange(p.id)} className="w-full">
          <div className="flex w-full items-center justify-between space-x-6 p-6">
            <div className="flex-1 truncate">
              <div className="flex items-center space-x-3">
                <h3 className="truncate text-sm font-medium text-gray-900">{p.title}</h3>
              </div>
            </div>
          </div>
          <div>
          </div>
        </button>
        </li>
      ))}
      <li className="col-span-1 divide-y divide-gray-200 rounded-lg bg-white shadow-sm">
        <button onClick={() => onChange()} className="w-full">
          <div className="flex w-full items-center justify-between space-x-6 p-6">
            <div className="flex-1 truncate">
              <div className="flex items-center space-x-3">
                <h3 className="truncate text-sm font-medium text-gray-900">Others</h3>
              </div>
            </div>
          </div>
          <div>
          </div>
          </button>
        </li>
    </ul>
  )
}
