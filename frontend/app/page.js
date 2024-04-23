"use client"
import React, { useState } from 'react';
import Navbar from '../components/navbar.js';

export default function Home() {
  return (
    <section
      className="rounded-lg bg-cover bg-no-repeat p-12 text-center relative object-cover"
      style={{backgroundImage: "url('/bg-website.png')", width: "100vw", height: "100vh"}}>
      <div
        className="absolute bottom-0 left-0 right-0 top-0 h-full w-full object-cover overflow-y-scroll bg-fixed"
        style={{backgroundColor: "rgba(0, 0, 0, 0.6)"}}>
        <div className="flex flex-col w-full items-center justify-center">
        <Navbar/>
        <img src="/Lemanspedia_Slogan-removebg.png" className='h-[300px] object-cover'/>
        {/*Start point What is WikiRace Card*/}
        <div className='rounded-[10px] my-2 w-[70%] h-fit relative border-2 border-white mx-4 mt-[50px]'>
            <div className='p-4 flex flex-col text-white'>
                <div className='flex items-start justify-center'>
                <h3 className='font-bold text-3xl text-white mx-[30px]'>What is <span className="underline">WikiRace</span>?</h3>
                </div>
                <div className='flex flex-row justify-center'>
                  <div className=''>
                    <img src={'\wikipedia-logo.png'} className='w-[400px] mt-[20px]'/>
                  </div>
                  <div className=' text-left w-[500px] ml-[25px]'>
                  <p className='font-inter text-l text-white mt-[6px] text-justify'> Wikiracing is a multiplayer virtual game themed around Wikipedia. The game measures the speed at which someone traverses links from one page to another. This game gained prominence among developers as it was once competed in TechOlympics and Yale Freshman Olympics.</p>
                  <p className='font-inter text-l text-white mt-[6px] text-justify'>The WikiRace processing on this website is done using the Go programming language. The website is capable of processing the entire route from the starting address to the destination Wikipedia address using the Iterative Deepening Search (IDS) and Breadth First Search (BFS) algorithms to complete the WikiRace game. The website can accept input in the form of algorithm type, initial article title, and destination article title. The output provided by the website includes the number of articles examined, the number of articles traversed, the route of article exploration (from the initial article to the destination article), and the search time (in ms).</p>
                  </div>
                </div>
              </div>
          </div>
          {/*end point What is WikiRace Card*/}
          {/* starpoint BFS IDS Algorithm */}
          <div className='justify-around w-full h-fit flex flex-row mb-[50px]'>
            {/* Algoritma BFS */}
          <div className='rounded-[10px] my-2 w-[32%] h-full relative border-2 border-white mx-4 mt-[50px] ml-[15%]'>
            <h2 className='font-bold text-2xl text-white'>Breadth First Search (BFS)</h2>
            <p className='font-inter text-l text-white mt-[6px] text-justify mx-[10px]'>Breadth-First Search (BFS) is a fundamental algorithm used in graph theory and computer science. It operates by exploring all the vertices (nodes) of a graph systematically, starting from a designated source vertex. The algorithm maintains a queue data structure to keep track of the vertices that need to be explored. At the beginning of the traversal, the source vertex is enqueued into the queue. Then, BFS iteratively dequeues vertices from the front of the queue and explores their adjacent vertices. This process continues until all vertices in the graph have been visited or until the desired condition is met. One of the key features of BFS is that it guarantees the shortest path from the source vertex to any other vertex in an unweighted graph. This property makes it particularly useful for tasks such as finding the shortest path between two nodes, determining connectivity, or exploring all reachable nodes in a graph. BFS is also commonly used in various applications such as network routing protocols, web crawlers, social network analysis, and shortest path algorithms in GPS navigation systems. Its simplicity and efficiency make it a versatile algorithm for solving a wide range of graph-related problems.</p>
          <div className='flex justify-center'>
          <img src='/BFS.gif' className='h-[250px] my-[70px]'/>
          </div>
          </div>
          {/* Algoritma IDS */}
          <div className='rounded-[10px] my-2 w-[32%] h-full relative border-2 border-white mx-4 mt-[50px] mr-[15%]'>
            <h2 className='font-bold text-2xl text-white'>Iterative Deepening Search (IDS)</h2>
            <p className='font-inter text-l text-white mt-[6px] text-justify mx-[10px]'>
Iterative Deepening Search (IDS) is a systematic search algorithm utilized for traversing trees or graphs. Operating with a strategy that amalgamates the benefits of Depth-First Search (DFS) and the completeness of Breadth-First Search (BFS), IDS iteratively conducts DFS with escalating depth limits until the target goal is attained. The algorithm initiates with an initial depth limit of 0 and proceeds to explore the search space. Initially, it performs DFS starting from the root node, restricting the depth of exploration to the current depth limit. Should the goal node be encountered within this limit, the solution is promptly returned. In cases where the goal node remains undiscovered but unexplored nodes exist at the current depth limit, the exploration continues. However, if DFS exhausts the exploration at the current depth limit without locating the goal node, the depth limit is incremented by 1, and the process is reiterated. This iterative approach persists until the goal node is successfully found or until all nodes within the search space have been explored. The utilization of iterative deepening in IDS allows for a balance between memory efficiency and search completeness, rendering it particularly suitable for scenarios where the depth of the solution is unknown or where memory constraints are pertinent.</p>
          <div className='flex justify-center'>
          <img src='/IDS.gif' className='h-[350px]'/>
          </div>
          </div>
          </div>
        </div>
      </div>
    </section>
  );
}
