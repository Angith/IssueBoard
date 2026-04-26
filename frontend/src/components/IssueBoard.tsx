'use client';

interface Issue {
  id: string;
  number: number;
  title: string;
  url: string;
  state: string;
}

interface Category {
  label: {
    name: string;
    color: string;
  };
  issues: Issue[];
}

interface Props {
  categories: Category[];
}

export default function IssueBoard({ categories }: Props) {
  return (
    <div className="flex gap-6 overflow-x-auto pb-8">
      {categories.map((category) => (
        <div key={category.label.name} className="min-w-[350px] max-w-[400px] flex-shrink-0">
          <div className="mb-4 flex items-center gap-2">
            <span
              className="h-3 w-3 rounded-full"
              style={{ backgroundColor: `#${category.label.color}` }}
            />
            <h2 className="text-xl font-bold">{category.label.name}</h2>
            <span className="rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600">
              {category.issues.length}
            </span>
          </div>

          <div className="space-y-4">
            {category.issues.map((issue) => (
              <a
                key={issue.id}
                href={issue.url}
                target="_blank"
                rel="noopener noreferrer"
                className="block rounded-lg border bg-white p-4 shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="mb-1 text-xs text-gray-500">#{issue.number}</div>
                <h3 className="font-medium text-gray-900 leading-tight">{issue.title}</h3>
              </a>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
