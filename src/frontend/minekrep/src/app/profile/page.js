'use client'

import NavBar from '@/components/NavBar';
import TeamProfile from '@/components/TeamProfile';

export default function ProfilePage() {
  const teamMembers = [
    {
      name: "Team Member 1",
      role: "Frontend Developer",
      nim: "13521XXX",
      image: "/images/team/avatar1.png",
      description: "Responsible for building the frontend using Next.js and implementing the React Flow visualization.",
      instagram: "https://instagram.com/team_member1",
      linkedin: "https://linkedin.com/in/team-member1",
      github: "https://github.com/team-member1"
    },
    {
      name: "Team Member 2",
      role: "Backend Developer",
      nim: "13521XXX",
      image: "/images/team/avatar2.png",
      description: "Worked on the Go backend implementation, including the BFS and DFS algorithms for recipe searches.",
      instagram: "https://instagram.com/team_member2",
      linkedin: "https://linkedin.com/in/team-member2",
      github: "https://github.com/team-member2"
    },
    {
      name: "Team Member 3",
      role: "Full Stack Developer",
      nim: "13521XXX",
      image: "/images/team/avatar3.png",
      description: "Integrated frontend and backend, implemented multithreading for recipe searches, and handled deployment.",
      instagram: "https://instagram.com/team_member3",
      linkedin: "https://linkedin.com/in/team-member3",
      github: "https://github.com/team-member3"
    }
  ];

  return (
    <div className="min-h-screen bg-gray-900"
      style={{ 
        backgroundImage: 'url("/images/minecraft_dark_background.jpg")', 
        backgroundSize: 'cover',
        fontFamily: '"Minecraft", sans-serif' 
      }}>
      <NavBar />
      
      <div className="container mx-auto px-4 py-12">
        <div className="bg-black bg-opacity-70 rounded-lg p-8 border-2 border-gray-700">
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold text-yellow-400 mb-4">
              Team Minekrep
            </h1>
            <p className="text-gray-300 max-w-2xl mx-auto">
              We are a group of students from Institut Teknologi Bandung (ITB) working on this project
              for the Algorithm Strategies course (IF2211). Our team has developed this application to demonstrate
              the use of BFS and DFS algorithms for finding recipes in Little Alchemy 2.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {teamMembers.map((member, index) => (
              <TeamProfile key={index} member={member} />
            ))}
          </div>
          
          <div className="mt-16 text-center border-t border-gray-700 pt-8">
            <h2 className="text-2xl font-bold text-green-400 mb-4">
              About the Project
            </h2>
            <p className="text-gray-300 mb-4">
              Little Alchemy 2 is a game where players combine elements to create new ones. Starting with the four basic
              elements (Earth, Fire, Water, Air), players can discover up to 720 different elements through various combinations.
            </p>
            <p className="text-gray-300 mb-4">
              Our application uses graph search algorithms (BFS and DFS) to find recipes for creating any element in the game. 
              BFS is great for finding the shortest path (fewest combinations), while DFS can explore all possible combinations.
            </p>
            <p className="text-gray-300">
              The project satisfies the requirements for the second major assignment (Tugas Besar 2) for the IF2211 Strategy Algorithm course
              at Institut Teknologi Bandung, academic year 2024/2025.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}