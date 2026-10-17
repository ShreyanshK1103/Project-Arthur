/**
 * ANSI → HTML via `anser`.
 * Chosen over ansi-to-react: smaller surface for a monospace log viewer.
 */
// Package uses `export =`; Vite interops the default at runtime.
// eslint-disable-next-line @typescript-eslint/no-require-imports
import AnserImport from "anser";

type AnserStatic = {
  ansiToHtml: (txt: string, options?: { use_classes?: boolean }) => string;
  escapeForHtml: (txt: string) => string;
};

const Anser = AnserImport as unknown as AnserStatic;

export function ansiToHtml(message: string): string {
  return Anser.ansiToHtml(Anser.escapeForHtml(message), {
    use_classes: false,
  });
}
