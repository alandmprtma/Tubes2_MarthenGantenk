import React, { useMemo } from "react";
import { ForceGraph2D } from "react-force-graph";

const linkColor = (link) => {
  if (link.type === 'special') {
    return 'red'; // Special links are red
  }
  return 'white'; // Default link color
};

const Graph = ({ path }) => {
  if (!Array.isArray(path)) {
    console.error('PathBox expects a prop "path" which is an array of paths.');
    return null;
  }

  const graphData = useMemo(() => {
    const nodes = new Set();
    const links = [];
    const startNode = path[0][0]; // Assume the start node is the first node of the first path
    const endNode = path[path.length - 1][path[path.length - 1].length - 1]; // Assume the end node is the last node of the last path

    path.forEach((singlePath) => {
      singlePath.forEach((node, index) => {
        nodes.add(node);
        if (index < singlePath.length - 1) {
          links.push({ source: node, target: singlePath[index + 1] });
        }
      });
    });

    return {
      nodes: Array.from(nodes).map((id) => ({ id, startNode: id === startNode, endNode: id === endNode })),
      links,
    };
  }, [path]);

  return <ForceGraph2D 
    graphData={graphData} 
    linkColor={link => linkColor(link)}
    nodeCanvasObject={(node, ctx, globalScale) => {
      const label = node.id;
      const fontSize = 12/globalScale;
      ctx.font = `${fontSize}px Sans-Serif`;
      const textWidth = ctx.measureText(label).width;
      const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2);

      ctx.fillStyle = node.startNode ? 'rgba(255, 0, 0, 0.8)' : node.endNode ? 'rgba(0, 255, 0, 0.8)' : 'rgba(255, 255, 255, 0.8)';
      ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, ...bckgDimensions);

      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillStyle = 'black';
      ctx.fillText(label, node.x, node.y);

      node.__bckgDimensions = bckgDimensions;
    }}
    minZoom={0.1}
  />;
};

export default Graph;
