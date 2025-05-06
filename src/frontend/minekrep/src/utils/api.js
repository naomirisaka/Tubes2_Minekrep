// API base URL - replace with your actual backend URL when deployed
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

/**
 * Search for element recipes
 * @param {Object} params - Search parameters
 * @param {string} params.algorithm - Search algorithm ('bfs', 'dfs', or 'bidirectional')
 * @param {string} params.targetElement - Target element to search for
 * @param {boolean} params.multipleRecipes - Whether to search for multiple recipes
 * @param {number} params.recipeCount - Number of recipes to search for (if multipleRecipes is true)
 * @returns {Promise<Object>} - Recipe search results
 */
export const searchRecipes = async ({
  algorithm = 'bfs',
  targetElement,
  multipleRecipes = false,
  recipeCount = 1
}) => {
  try {
    const response = await fetch(`${API_BASE_URL}/search`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        algorithm,
        targetElement,
        multipleRecipes,
        recipeCount,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to search for recipes');
    }

    return await response.json();
  } catch (error) {
    console.error('Error searching recipes:', error);
    throw error;
  }
};

/**
 * Get all available elements
 * @returns {Promise<Array>} - List of all available elements
 */
export const getAllElements = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/elements`);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to fetch elements');
    }

    return await response.json();
  } catch (error) {
    console.error('Error fetching elements:', error);
    throw error;
  }
};

/**
 * Get basic elements
 * @returns {Promise<Array>} - List of basic elements
 */
export const getBasicElements = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/elements/basic`);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to fetch basic elements');
    }

    return await response.json();
  } catch (error) {
    console.error('Error fetching basic elements:', error);
    throw error;
  }
};