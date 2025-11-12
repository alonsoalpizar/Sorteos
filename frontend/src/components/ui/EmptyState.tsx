import { ReactNode } from 'react';
import { cn } from '@/lib/utils';
import { GradientButton } from './GradientButton';

interface EmptyStateProps {
  icon: ReactNode;
  title: string;
  description: string;
  action?: {
    label: string;
    onClick: () => void;
  };
  className?: string;
}

export function EmptyState({ icon, title, description, action, className }: EmptyStateProps) {
  return (
    <div className={cn("text-center py-16 px-4 animate-fade-in", className)}>
      <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20 rounded-full mb-6 text-primary-600 dark:text-primary-400">
        {icon}
      </div>
      <h3 className="text-2xl font-bold text-slate-900 dark:text-white mb-3">
        {title}
      </h3>
      <p className="text-slate-600 dark:text-slate-400 mb-8 max-w-md mx-auto leading-relaxed">
        {description}
      </p>
      {action && (
        <GradientButton onClick={action.onClick} variant="primary">
          {action.label}
        </GradientButton>
      )}
    </div>
  );
}
