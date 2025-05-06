'use client'

import { useState, useEffect, useCallback } from 'react';
import ReactFlow, {
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  MarkerType,
} from 'reactflow';
import 'reactflow/dist/style.css';

// Custom node renderer for element nodes
const ElementNode = ({ data }) => {
  return (
    <div className={`
      p-2 rounded-md min-w-32 text-center font-semibold
      ${data.isBasicElement 
        ? 'bg-green-700 border-2 border-green-500' 
        : data.isCombineNode 
          ? 'bg-gray-600 border-2 border-gray-500' 
          : 'bg-blue-700 border-2 border-blue-500'}
    `}>
      <div className="flex items-center justify-center">
        {data.icon && (
          <img src={data.icon} alt={data.label} className="w-6 h-6 mr-2" />
        )}
        <div>{data.label}</div>
      </div>
    </div>
  );
};

// Node types mapping
const nodeTypes = {
  elementNode: ElementNode,
};

const RecipeVisualizer = ({ recipes }) => {
  const [recipeIndex, setRecipeIndex] = useState(0);
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [isLoading, setIsLoading] = useState(true);

  // Helper to get element icon (mock function, real implementation would fetch icons)
  const getElementIcon = (element) => {
    const basicElements = {
      'earth': '/images/elements/earth.png',
      'fire': '/images/elements/fire.png',
      'water': '/images/elements/water.png',
      'air': '/images/elements/air.png'
    };
    
    return basicElements[element.toLowerCase()] || '/images/elements/generic.png';
  };
  
  // Create the recipe tree visualization when recipes or recipeIndex change
  useEffect(() => {
    if (!recipes || recipes.length === 0 || recipeIndex >= recipes.length) {
      setNodes([]);
      setEdges([]);
      setIsLoading(false);
      return;
    }
    
    setIsLoading(true);
    
    const currentRecipe = recipes[recipeIndex];
    const newNodes = [];
    const newEdges = [];
    let nodeId = 0;

    // Recursive function to build the tree
    const buildTree = (element, level, parentId = null, position = { x: 0, y: 0 }) => {
      const currentId = `node_${nodeId++}`;
      const isBasicElement = ['earth', 'fire', 'water', 'air'].includes(element.name.toLowerCase());
      const isCombineNode = element.name === '+';
      
      // Create the node
      newNodes.push({
        id: currentId,
        type: 'elementNode',
        position,
        data: {
          label: element.name,
          icon: getElementIcon(element.name),
          isBasicElement,
          isCombineNode
        }
      });
      
      // Create edge to parent if it exists
      if (parentId) {
        newEdges.push({
          id: `edge_${parentId}_${currentId}`,
          source: parentId,
          target: currentId,
          type: 'step',
          markerEnd: {
            type: MarkerType.ArrowClosed,
            width: 15,
            height: 15,
          },
          style: { stroke: '#aaa' }
        });
      }
      
      // Process children if they exist
      if (element.children && element.children.length > 0) {
        // Calculate horizontal spacing for children
        const childSpacing = 150; // horizontal space between children
        const startX = position.x - ((element.children.length - 1) * childSpacing) / 2;
        
        element.children.forEach((child, index) => {
          const childPos = {
            x: startX + (index * childSpacing),
            y: position.y + 100 // vertical spacing
          };
          
          buildTree(child, level + 1, currentId, childPos);
        });
      }
      
      return currentId;
    };
    
    // Start building from the final element (the root of our tree)
    const rootPosition = { x: 300, y: 50 }; // Center the tree
    buildTree(currentRecipe, 0, null, rootPosition);
    
    setNodes(newNodes);
    setEdges(newEdges);
    setIsLoading(false);
  }, [recipes, recipeIndex, setNodes, setEdges]);
  
  // Layout for the diagram
  const fitView = useCallback(() => {
    if (reactFlowInstance) {
      reactFlowInstance.fitView({ padding: 0.2 });
    }
  }, []);
  
  const [reactFlowInstance, setReactFlowInstance] = useState(null);
  
  return (
    <div className="flex flex-col h-full">
      {recipes && recipes.length > 1 && (
        <div className="flex items-center justify-between mb-4">
          <button
            onClick={() => setRecipeIndex(prev => Math.max(0, prev - 1))}
            disabled={recipeIndex === 0}
            className={`
              px-3 py-1 rounded-md text-sm
              ${recipeIndex === 0 
                ? 'bg-gray-700 text-gray-400 cursor-not-allowed' 
                : 'bg-blue-700 hover:bg-blue-600 text-white'}
            `}
          >
            Previous Recipe
          </button>
          
          <div className="text-gray-300 text-sm">
            Recipe {recipeIndex + 1} of {recipes.length}
          </div>
          
          <button
            onClick={() => setRecipeIndex(prev => Math.min(recipes.length - 1, prev + 1))}
            disabled={recipeIndex === recipes.length - 1}
            className={`
              px-3 py-1 rounded-md text-sm
              ${recipeIndex === recipes.length - 1 
                ? 'bg-gray-700 text-gray-400 cursor-not-allowed' 
                : 'bg-blue-700 hover:bg-blue-600 text-white'}
            `}
          >
            Next Recipe
          </button>
        </div>
      )}
      
      {isLoading ? (
        <div className="flex justify-center items-center h-80">
          <div className="text-gray-400">Building recipe visualization...</div>
        </div>
      ) : recipes && recipes.length > 0 ? (
        <div style={{ height: '60vh' }} className="border border-gray-700 rounded-md overflow-hidden">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            nodeTypes={nodeTypes}
            onInit={setReactFlowInstance}
            fitView
            attributionPosition="bottom-right"
          >
            <Controls />
            <Background color="#444" gap={16} />
          </ReactFlow>
        </div>
      ) : (
        <div className="flex justify-center items-center h-40 text-gray-400">
          No recipes found for this element. Try another search.
        </div>
      )}
    </div>
  );
};

export default RecipeVisualizer;