'use client'

import { useState } from 'react';
import MinecraftButton from './MinecraftButton';

const SearchForm = ({ 
  targetElement, 
  setTargetElement, 
  algorithm,
  setAlgorithm,
  shortestPath, 
  setShortestPath, 
  recipeCount, 
  setRecipeCount, 
  handleSearch, 
  loading 
}) => {
  return (
    <form onSubmit={handleSearch} className="space-y-4">
      {/* Element Input */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          Target Element
        </label>
        <input
          type="text"
          value={targetElement}
          onChange={(e) => setTargetElement(e.target.value)}
          placeholder="Enter element name (e.g., Brick, Metal, Human)"
          className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md 
                    focus:ring-2 focus:ring-green-500 focus:border-transparent
                    placeholder-gray-400 text-white"
          required
        />
        <p className="text-xs text-gray-400 mt-1">
          Search for any of the 720 elements available in Little Alchemy 2
        </p>
      </div>
      
      {/* Toggle for Shortest Path or Multiple Recipes */}
      <div className="pt-2">
        <label className="block text-sm font-medium text-gray-300 mb-3">
          Search Type
        </label>
        
        <div className="flex gap-4">
          <div 
            className={`
              flex-1 p-3 rounded-md cursor-pointer transition-all
              ${shortestPath 
                ? 'bg-green-900 border-2 border-green-600' 
                : 'bg-gray-700 border-2 border-gray-600'}
            `}
            onClick={() => setShortestPath(true)}
          >
            <div className="flex items-center">
              <div className={`
                w-4 h-4 rounded-full mr-2
                ${shortestPath ? 'bg-green-400' : 'bg-gray-500'}
                flex items-center justify-center
              `}>
                {shortestPath && (
                  <div className="w-2 h-2 bg-white rounded-full"></div>
                )}
              </div>
              <div>
                <h4 className="font-medium text-sm">One Recipe</h4>
                <p className="text-xs text-gray-400">Find one way to create</p>
              </div>
            </div>
          </div>
          
          <div 
            className={`
              flex-1 p-3 rounded-md cursor-pointer transition-all
              ${!shortestPath 
                ? 'bg-green-900 border-2 border-green-600' 
                : 'bg-gray-700 border-2 border-gray-600'}
            `}
            onClick={() => setShortestPath(false)}
          >
            <div className="flex items-center">
              <div className={`
                w-4 h-4 rounded-full mr-2
                ${!shortestPath ? 'bg-green-400' : 'bg-gray-500'}
                flex items-center justify-center
              `}>
                {!shortestPath && (
                  <div className="w-2 h-2 bg-white rounded-full"></div>
                )}
              </div>
              <div>
                <h4 className="font-medium text-sm">Multiple Recipes</h4>
                <p className="text-xs text-gray-400">Find different ways to create</p>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      {/* Recipe Count Slider (only visible when Multiple Recipes is selected) */}
      {!shortestPath && (
        <div className="pt-1">
          <label className="block text-sm font-medium text-gray-300 mb-1">
            Max Recipe Count: {recipeCount}
          </label>
          <input
            type="range"
            min="1"
            max="10"
            value={recipeCount}
            onChange={(e) => setRecipeCount(parseInt(e.target.value))}
            className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-green-600"
          />
          <div className="flex justify-between text-xs text-gray-400">
            <span>1</span>
            <span>5</span>
            <span>10</span>
          </div>
        </div>
      )}
      
      {/* Submit Button */}
      <div className="pt-2">
        <MinecraftButton
          text={loading ? "Searching..." : "Find Recipes"}
          variant="primary"
          className="w-full"
          disabled={loading}
          type="submit"
        />
      </div>
    </form>
  );
};

export default SearchForm;