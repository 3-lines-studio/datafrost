declare module "prismjs/components/prism-core" {
  export function highlight(
    code: string,
    grammar: any,
    language: string,
  ): string;
  export const languages: Record<string, any>;
}

declare module "sql-formatter" {
  export function format(sql: string, options?: any): string;
}
