import { Tag } from "antd";

// Renders a string-list JSON field as chips, used wherever a contact's tags show.
export function TagChips({ tags }: { tags: string[] }) {
  if (tags.length === 0) return "-";
  return (
    <span className="inline-flex flex-wrap gap-y-1">
      {tags.map((tag) => (
        <Tag key={tag} variant="filled" color="default">
          {tag}
        </Tag>
      ))}
    </span>
  );
}
