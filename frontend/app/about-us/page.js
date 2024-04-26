import React from 'react';
import Navbar from '../../components/navbar.js';
// import aland from '/aland.jpg';
// import qika from '/qika.jpg';
// import Ikhwan from '/ikhwan.jpg';

export default function Aboutus() {
    const members = [
        {
          name: 'Aland Mulia Pratama',
          nim: '13522124',
          role: 'Frontend & Backend Developer',
          email: '13522124@std.stei.itb.ac.id',
          image: '/aland.jpg',
        },
        {
          name: 'Rizqika Mulia Pratama',
          nim: '13522126',
          role: 'IDS Algorithm & Backend Developer',
          email: '13522126@std.stei.itb.ac.id',
          image: '/qika.jpg',
        },
        {
          name: 'Ikhwan Al Hakim',
          nim: '13522147',
          role: 'BFS Algorithm & Backend Developer',
          email: '13522147@std.stei.itb.ac.id',
          image: '/ikhwan.jpg',
        },
      ];

  return (
      <div className="flex flex-col w-full items-center justify-center">
      <Navbar/>
      <h1 className='font-bold text-3xl text-white mt-3'>LEMANSPEDIA CONTRIBUTORS</h1>
      <h3 className=' text-xl text-white text-center mx-[50px]'>The main objective from the second major assignment of Algorithm Strategies courses is to make a WikiRace programs with Breadth First Search (BFS) and Iterative Deepening Search (IDS) Algorithms.</h3>
      <article className='flex w-[80%] gap-x-4 justify-center mt-10'>
        {members.map((items) => {
          return (
          <div className='rounded-[10px] my-2 w-[300px] h-[350px] relative border-2 border-white mx-4'>
            <div className='flex flex-col items-center mt-7'>
              <img src={items.image} className='h-[150px] w-[150px] rounded-full object-cover z-10'/>
              <p className='font-bold mt-7'>{items.name}</p>
              <p className='font'>{items.nim}</p>
              <p className='italic'>{items.role}</p>
              <p className='font'>({items.email})</p>
            </div>
          </div>
          )
        })}
      </article>
      </div>
    );
}
