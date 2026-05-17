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
  error?: string;
}

export default function IssueBoard({ categories, error }: Props) {
  if (error) {
    return (
      <div className="p-8 border border-red-900/30 bg-red-950/20 rounded-xl text-red-400">
        <h3 className="text-xl font-semibold mb-2 text-red-300">Something went wrong</h3>
        <p>{error}</p>
        {error.includes('rate limit') && (
          <p className="mt-4 text-sm">
            GitHub unauthenticated API has a limit of 60 requests per hour. 
            Please wait a while before refreshing again.
          </p>
        )}
      </div>
    );
  }

  if (categories.length === 0) {
    return (
      <div className="p-8 text-center text-zinc-500 border border-dashed border-zinc-800 bg-zinc-900/20 rounded-xl">
        No issues found for this repository.
      </div>
    );
  }

  return (
    <div className="flex gap-6 overflow-x-auto pb-8">
      {categories.map((category) => (
        <div key={category.label.name} className="min-w-[350px] max-w-[400px] flex-shrink-0">
          <div className="mb-4 flex items-center gap-2">
            <span
              className="h-3 w-3 rounded-full shadow-sm"
              style={{ backgroundColor: `#${category.label.color}` }}
            />
            <h2 className="text-sm font-medium text-zinc-200">{category.label.name}</h2>
            <span className="rounded-full bg-zinc-800 px-2 py-0.5 text-xs font-medium text-zinc-400">
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
                className="block rounded-xl border border-zinc-800/60 bg-zinc-900/50 p-4 transition-all hover:border-zinc-700 hover:bg-zinc-900"
              >
                <div className="mb-2 text-xs font-medium text-zinc-500">#{issue.number}</div>
                <h3 className="text-sm font-medium text-zinc-200 leading-snug">{issue.title}</h3>
              </a>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
