import { Todo, updateTodo } from "../api/todo";

const StatusMap: { [key: number]: string } = {
    0: "Not Done",
    1: "In Progress",
    2: "Done"
}

const PriorityMap: { [key: number]: string } = {
    0: "Low",
    1: "Medium",
    2: "High"
}

export default function TodoTable({todos}: {todos: Todo[]}) {
    function handleClick(id: number) {
        updateTodo(id);
        window.location.reload();
    }
    return (
        <div className="mt-8 flow-root">
        <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle">
            <table className="relative min-w-full divide-y divide-gray-300">
              <thead>
                <tr>
                  <th
                    scope="col"
                    className="py-3.5 pr-3 pl-4 text-left text-sm font-semibold text-gray-900 sm:pl-6 lg:pl-8"
                  >
                    Title
                  </th>
                  <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                    Priority
                  </th>
                  <th scope="col" className="py-3.5 pr-4 pl-3 sm:pr-6 lg:pr-8">
                    <span className="sr-only">Edit</span>
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 bg-white">
                {todos && todos.map((td) => (
                  <tr key={td.id}>
                    <td className="py-4 pr-3 pl-4 text-sm font-medium whitespace-nowrap text-gray-900 sm:pl-6 lg:pl-8">
                      {td.title}
                    </td>
                    <td className="px-3 py-4 text-sm whitespace-nowrap text-gray-500">{PriorityMap[td.priority]}</td>
                    <td className="py-4 pr-4 pl-3 text-right text-sm font-medium whitespace-nowrap sm:pr-6 lg:pr-8">
                      <button
                            type="button"
                            onClick={(e) => {handleClick(td.id)}}
                            aria-pressed={td.status === 2}
                            aria-label={td.status!=2 ? "Mark as not done" : "Mark as done"}
                            className={`inline-flex items-center justify-center w-10 h-10 rounded-full transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-400
                                ${td.status === 2 ? "bg-green-500 border-transparent shadow-sm" : "bg-white border-2 border-gray-200"}`}
                            >
                            {td.status === 2 ? (
                                <svg className="w-5 h-5 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                <path d="M20 6L9 17l-5-5" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                                </svg>
                            ) : (
                                <svg className="w-5 h-5 text-gray-500" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                <circle cx="12" cy="12" r="9" strokeWidth="2"></circle>
                                </svg>
                            )}
                            </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    )

}