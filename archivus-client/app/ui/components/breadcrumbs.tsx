import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";


export default function PathBreadcrumbs({parts}:{parts:string[]}) {
  return (
    <Breadcrumb>
    <BreadcrumbList>
    <BreadcrumbItem key={-1}>
            <BreadcrumbLink href='/'>
              Home
            </BreadcrumbLink>
    </BreadcrumbItem>
    {parts.length > 0 && <BreadcrumbSeparator />}
    {
      parts.map((part, index) => {
        const href = '/?folder=' + parts.slice(0, index + 1).join('/');
        return (
            <div key={index} className="flex items-center">
          <BreadcrumbItem key={index}>
            <BreadcrumbLink href={href}>
              {part}
            </BreadcrumbLink>
          </BreadcrumbItem>
        {index < parts.length - 1 && <BreadcrumbSeparator />}
            </div>
        )
      })
    }
    
    </BreadcrumbList>
    
    </Breadcrumb>
  )
}