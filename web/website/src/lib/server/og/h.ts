// Minimal element builder for satori, so templates read like markup instead of nested objects.

export type Node = { type: string; props: Record<string, unknown> } | string | null | false;
type Style = Record<string, string | number>;

export function h(type: string, style: Style, ...children: Node[]): Node {
  const kids = children.filter((c) => c !== null && c !== false);
  return {
    type,
    props: {
      style: type === 'div' ? { display: 'flex', ...style } : style,
      children: kids.length === 1 ? kids[0] : kids,
    },
  };
}

/** An element with extra attributes (used for SVG). */
export function el(type: string, attrs: Record<string, unknown>, ...children: Node[]): Node {
  return { type, props: { ...attrs, children } };
}
