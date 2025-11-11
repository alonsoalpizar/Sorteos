interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  text?: string;
}

export function LoadingSpinner({ size = 'md', text }: LoadingSpinnerProps) {
  const sizeClasses = {
    sm: 'w-5 h-5',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  };

  return (
    <div className="flex flex-col items-center justify-center py-12">
      <div
        className={`${sizeClasses[size]} border-4 border-slate-200 dark:border-slate-700 border-t-blue-600 rounded-full animate-spin`}
      />
      {text && (
        <p className="mt-4 text-sm text-slate-600 dark:text-slate-400">
          {text}
        </p>
      )}
    </div>
  );
}
