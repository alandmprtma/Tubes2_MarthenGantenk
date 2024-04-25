import React from "react";

const PathBox = ({ path }) => {
  // Ensure the 'path' prop is an array
  if (!Array.isArray(path)) {
    console.error('PathBox expects a prop "path" which is an array of strings.');
    return null;
  }

  return (
    <div className="p-8">
      <div className=" text-gray-600 flex flex-wrap items-center gap-8 justify-center">
       {path.map((paths, index) => (
        <p key={index} className="text-white">
        {paths.map((node, nodeIndex) => (
          <div
            key={nodeIndex}
            className="box-border p-4 border-2 border-gray-200 rounded-md shadow-sm flex items-center justify-between"
            style={{ backgroundColor: 'rgba(255, 255, 255, 0.1)' }}
          >
          <React.Fragment key={nodeIndex}>
            <a
              className="text-white hover:text-[#ADD8E6]"
              href={`https://en.wikipedia.org/wiki/${encodeURIComponent(node)}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              {node}
            </a>
            {/* Render the arrow if this is not the last item */}
            <span className="text-white mx-2">→</span>
          </React.Fragment>
          </div>
        ))}
         </p>
         ))}
      </div>
    </div>
  );
};

export default PathBox;