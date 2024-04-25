"use client"
import React, {useCallback, useMemo, useState } from "react";
import { ForceGraph2D } from "react-force-graph";

const linkColor = (link) => {
    // You can determine the color based on the properties of the link
    // For example, if your link objects have a 'type' property:
    if (link.type === 'special') {
      return 'red'; // Special links are red
    }
    return 'white'; // Default link color is blue
  };

 

const Graph = ({ path }) => {
  const [zoom, setZoom] = useState(1); // State to keep track of zoom level
  const handleZoom = useCallback((newZoom) => {
    setZoom(newZoom);
  }, []);

  // Ensure the 'path' prop is an array
  if (!Array.isArray(path)) {
    console.error('PathBox expects a prop "path" which is an array of paths.');
    return null;
  }

  // Transform path data into graph data for visualization
  const graphData = useMemo(() => {
    const nodes = new Set();
    const links = [];



    path.forEach((singlePath) => {
      singlePath.forEach((node, index) => {
        nodes.add(node); // Add the node if it's not already in the set
        if (index < singlePath.length - 1) {
          links.push({ source: node, target: singlePath[index + 1] }); // Add the link between the current node and the next one
        }
      });
    });

    return {
      nodes: Array.from(nodes).map((id) => ({ id })),
      links,
    };
  }, [path]);

  return <ForceGraph2D 
  graphData={graphData} 
  onZoom={handleZoom}
  linkColor={link => linkColor(link)}
  nodeCanvasObject={(node, ctx, globalScale) => {
    const label = node.id;
    const fontSize = 12/globalScale; // Adjust font size to scale with zoom
    ctx.font = `${fontSize}px Sans-Serif`;
    const textWidth = ctx.measureText(label).width;
    const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2); // some padding

    ctx.fillStyle = 'rgba(255, 255, 255, 0.8)'; // Background color
    ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, ...bckgDimensions);

    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillStyle = 'black'; // Text color
    ctx.fillText(label, node.x, node.y);

    node.__bckgDimensions = bckgDimensions; // to re-use in nodePointerAreaPaint
  }}
  nodePointerAreaPaint={(node, color, ctx) => {
    ctx.fillStyle = color;
    const bckgDimensions = node.__bckgDimensions;
    bckgDimensions && ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, ...bckgDimensions);
  }}
  maxZoom={2}  
  minZoom={0.1}
/>;
};

export default Graph;