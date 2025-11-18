import { ReactNode } from "react";
import { clsx } from "clsx";

interface TableProps {
  children: ReactNode;
  className?: string;
}

interface TableHeaderProps {
  children: ReactNode;
  className?: string;
}

interface TableRowProps {
  children: ReactNode;
  className?: string;
  onClick?: () => void;
}

interface TableHeadProps {
  children: ReactNode;
  className?: string;
}

interface TableCellProps {
  children: ReactNode;
  className?: string;
}

export function Table({ children, className }: TableProps) {
  return (
    <div className="w-full overflow-x-auto">
      <table className={clsx("w-full border-collapse", className)}>
        {children}
      </table>
    </div>
  );
}

export function TableHeader({ children, className }: TableHeaderProps) {
  return (
    <thead className={clsx("bg-slate-50 border-b border-slate-200", className)}>
      {children}
    </thead>
  );
}

export function TableBody({ children, className }: TableProps) {
  return <tbody className={className}>{children}</tbody>;
}

export function TableRow({ children, className, onClick }: TableRowProps) {
  return (
    <tr
      className={clsx(
        "border-b border-slate-200 last:border-0 transition-colors",
        onClick && "cursor-pointer hover:bg-slate-50",
        className
      )}
      onClick={onClick}
    >
      {children}
    </tr>
  );
}

export function TableHead({ children, className }: TableHeadProps) {
  return (
    <th
      className={clsx(
        "px-4 py-3 text-left text-xs font-medium text-slate-700 uppercase tracking-wider",
        className
      )}
    >
      {children}
    </th>
  );
}

export function TableCell({ children, className }: TableCellProps) {
  return (
    <td className={clsx("px-4 py-3 text-sm text-slate-900", className)}>
      {children}
    </td>
  );
}
