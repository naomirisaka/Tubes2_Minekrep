'use client'

import { useState, useEffect } from 'react';
import Link from 'next/link';
import NavBar from '@/components/NavBar';
import AlgorithmSelector from '@/components/AlgorithmSelector';
import SearchForm from '@/components/SearchForm';
import RecipeVisualizer from '@/components/RecipeVisualizer';
import { searchRecipes } from '@/utils/api';

export default function SearchPage() {
  const [algorithm, setAlgorithm] = useState('bfs');
  const [targetElement, setTargetElement] = useState('');
  const [recipeCount, setRecipeCount] = useState(1);
  const [shortestPath, setShortestPath] = useState(true);
  const [searchResults, setSearchResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [metrics, setMetrics] = useState({ time: 0, nodesVisited: 0 });
  
  const handleSearch = async (e) => {
    e.preventDefault();
    
    if (!targetElement) {
      setError('Please enter an element name to search for');
      return;
    }
    
    setLoading(true);
    setError(null);
    setSearchResults(null);
    
    try {
      const results = await searchRecipes({
        algorithm,
        targetElement,
        multipleRecipes: !shortestPath,
        recipeCount: shortestPath ? 1 : recipeCount
      });
      
      setSearchResults(results.recipes);
      setMetrics({
        time: results.metrics.time,
        nodesVisited: results.metrics.nodesVisited
      });
    } catch (err) {
      console.error('Search error:', err);
      setError(err.message || 'Failed to search for recipes. Please try again.');
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <div className="min-h-screen bg-gray-900 text-white"
      style={{ 
        backgroundImage: 'url("/images/minecraft_dark_background.jpg")', 
        backgroundSize: 'cover',
        fontFamily: '"Minecraft", sans-serif' 
      }}>
      <NavBar />
      
      <div className="container mx-auto px-4 py-8">
        <div className="bg-black bg-opacity-80 rounded-lg p-6 border-2 border-gray-700 mb-8">
          <h1 className="text-3xl font-bold mb-6 text-center text-yellow-400">
            Recipe Finder
          </h1>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div className="bg-gray-800 bg-opacity-80 p-4 rounded-lg border border-gray-700">
              <h2 className="text-xl font-semibold mb-4 text-green-400">Algorithm</h2>
              <AlgorithmSelector 
                algorithm={algorithm} 
                setAlgorithm={setAlgorithm} 
              />
            </div>
            
            <div className="bg-gray-800 bg-opacity-80 p-4 rounded-lg border border-gray-700">
              <h2 className="text-xl font-semibold mb-4 text-green-400">Search Options</h2>
              <SearchForm 
                targetElement={targetElement}
                setTargetElement={setTargetElement}
                shortestPath={shortestPath}
                setShortestPath={setShortestPath}
                recipeCount={recipeCount}
                setRecipeCount={setRecipeCount}
                handleSearch={handleSearch}
                loading={loading}
              />
            </div>
          </div>
          
          {error && (
            <div className="bg-red-900 bg-opacity-80 p-4 rounded-lg border border-red-700 mb-6 text-white">
              {error}
            </div>
          )}
          
          {loading && (
            <div className="flex justify-center items-center p-12">
              <div className="animate-bounce text-xl text-green-400">
                Searching recipes...
              </div>
            </div>
          )}
          
          {searchResults && !loading && (
            <div className="bg-gray-800 bg-opacity-80 p-4 rounded-lg border border-gray-700 mt-6">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-semibold text-green-400">Recipe Results</h2>
                <div className="text-sm text-gray-400">
                  <div>Time: {metrics.time.toFixed(2)}ms</div>
                  <div>Nodes visited: {metrics.nodesVisited}</div>
                </div>
              </div>
              
              <RecipeVisualizer recipes={searchResults} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}