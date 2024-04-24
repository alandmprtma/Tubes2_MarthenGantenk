"use client"
import React, { useState } from 'react';
import Navbar from '../../components/navbar.js';
import { css } from "@emotion/react";
import { BeatLoader } from "react-spinners";
import axios from 'axios';
// import Particles from 'react-particles-js';


const override = css`
  /* Definisikan gaya khusus di sini */
`;

export default function Wikirace() {
  function delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
  const [awal, setAwal] = useState('');
  const [akhir, setAkhir] = useState('');
  const [inputsFilled, setInputsFilled] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [loading, setLoading] = useState(false);
  const [openAwal, setOpenAwal] = useState(false);
  const [openAkhir, setOpenAkhir] = useState(false);
  // State untuk autocomplete
  const [resultAwal, setResultAwal]= useState([])
  const [resultAkhir, setResultAkhir]= useState([])
  const [activeAlgorithm, setActiveAlgorithm] = useState(''); // default to BFS
  const [activeSolution, setActiveSolution] = useState(''); // default to Single Solution

  const handleAlgorithmClick = (algorithm) => {
    console.log("Algorithm", algorithm);
    setActiveAlgorithm(algorithm);
  };

  const handleSolutionClick = (solution) => {
    setActiveSolution(solution);
  };

  const handleChangeAwal = (event) => {
    setOpenAwal(true);
    handleQueryAwal();
    setAwal(event.target.value);
    setInputsFilled(event.target.value !== '' || akhir !== '');
  };

  const handleChangeAkhir = (event) => {
    setOpenAkhir(true);
    handleQueryAkhir();
    setAkhir(event.target.value);
    setInputsFilled(awal !== '' || event.target.value !== '');
  };

  const handleSubmit = async (event) => {
    setSubmitted(false);
    setLoading(true); // Set loading to true
    event.preventDefault();
    try {
      // Lakukan proses pengiriman data
      await delay(2000); // Misalnya, panggil API
      // Setelah selesai, atur state loading menjadi false
      setLoading(false);
      // Atur submitted menjadi true setelah proses selesai
      setSubmitted(true);
    } catch (error) {
      console.error(error);
      setLoading(false); // Jika terjadi kesalahan, set loading menjadi false
    }
  };

    const handleQueryAwal = async () => {
      const value = awal.trim();

      // if (!value) {
      //     console.error("No query provided");
      //     return;
      // }

      try {
          const response = await axios.get(
              `http://localhost:8080/api/wikipedia?query=${encodeURIComponent(value)}`
          );

          console.log("Response data:", response.data); // Check response structure

          // Assuming response.data is directly an array of results as per your backend code
          if (Array.isArray(response.data)) {
              const results = response.data.map(item => ({
                  title: item.title,
                  thumbnail: item.thumbnail || "", // Handle missing thumbnail
              }));

              console.log("Formatted results:", results);
              setResultAwal(results); // Assuming setResultAwal is your state setter
          } else {
              console.error('Error fetching data: Invalid response format');
          }
      } catch (error) {
          console.error('Error fetching data:', error);
      }
  };

    
  const handleQueryAkhir = async () => {
    const value = akhir.trim();

    // if (!value) {
    //     console.error("No query provided");
    //     return;
    // }

    try {
        const response = await axios.get(
            `http://localhost:8080/api/wikipedia?query=${encodeURIComponent(value)}`
        );

        console.log("Response data:", response.data); // Check response structure

        // Assuming response.data is directly an array of results as per your backend code
        if (Array.isArray(response.data)) {
            const results = response.data.map(item => ({
                title: item.title,
                thumbnail: item.thumbnail || "", // Handle missing thumbnail
            }));

            console.log("Formatted results:", results);
            setResultAkhir(results); // Assuming setResultAwal is your state setter
        } else {
            console.error('Error fetching data: Invalid response format');
        }
    } catch (error) {
        console.error('Error fetching data:', error);
    }
};

 // Button style base
 const baseStyle = "mx-4 rounded border-2 px-7 pb-[8px] pt-[10px] text-sm font-medium uppercase leading-normal transition duration-150 ease-in-out focus:outline-none focus:ring-0";
 // Common dynamic style
 const dynamicStyle = (isActive) => 
   `${baseStyle} ${isActive ? 'border-neutral-100 text-neutral-100 bg-neutral-500 bg-opacity-50' : 'border-neutral-50 text-neutral-50 hover:border-neutral-100 hover:bg-neutral-500 hover:bg-opacity-10 hover:text-neutral-100'}`;

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
          <div className="text-white mt-[25px]">
            <h2 className="mb-4 text-2xl font-semibold">Find the shortest paths from</h2>
            <form onSubmit={handleSubmit}>
              <label className=' mb-[20px] flex flex-row'>
                <div flex flex-col>
                <input
                  className='w-[500px] h-[60px] font-inter rounded-[10px] border-2 border-white mr-2'
                  style={{ backgroundColor: 'rgba(255, 255, 255, 0)' }}
                  type="text"
                  value={awal}
                  onChange={handleChangeAwal}
                  placeholder="Masukkan alamat awal"
                />
                {resultAwal && openAwal && (
                  <div className='text-white text-l flex flex-col items-center h-[200px] overflow-y-scroll absolute'>
                    {resultAwal.map((item, index) => {
                      return (
                        <div className='border-[1px] w-[475px] border-white flex flex-row items-center p-2 mb-2 rounded-lg bg-[#333]'
                          key={index}
                          onClick={() => {
                            setAwal(item.title)
                            setOpenAwal((prev) =>!prev)
                          }}
                        >
                          {item.thumbnail && (
                            <img src={item.thumbnail} alt={item.title} className='w-10 h-10 mr-2 rounded-full'/>
                          )}
                          <div>
                            <p className='font-bold'>{item.title}</p>
                          </div>
                        </div>
                      )
                    })}
                  </div>
                )}
                </div>
                {inputsFilled ? (
                  <img src="arrow-right-arrow-left-solid.svg" alt="to" className="w-[25px] mx-3 mt-4 mb-4 h-[25px] top-0" />
                ) : (
                  <h2 className="mx-3 mt-4 mb-4 text-2xl font-semibold h-[60px]">to</h2>
                )}
                <div flex flex-col>
                <input
                  className='w-[500px] h-[60px] font-inter rounded-[10px] border-2 border-white mr-2'
                  style={{ backgroundColor: 'rgba(255, 255, 255, 0)' }}
                  type="text"
                  value={akhir}
                  onChange={handleChangeAkhir}
                  placeholder="Masukkan alamat akhir"
                />
                {resultAkhir && openAkhir && (
                  <div className='text-white text-l flex flex-col items-center h-[200px] overflow-y-scroll absolute'>
                    {resultAkhir.map((item, index) => {
                      return (
                        <div className='border-[1px] w-[475px] border-white flex flex-row items-center p-2 mb-2 rounded-lg bg-[#333]'
                          key={index}
                          onClick={() => {
                            setAkhir(item.title)
                            setOpenAkhir((prev) =>!prev)
                          }}
                        >
                          {item.thumbnail && (
                            <img src={item.thumbnail} alt={item.title} className='w-10 h-10 mr-2 rounded-full'/>
                          )}
                          <div>
                            <p className='font-bold'>{item.title}</p>
                          </div>
                        </div>
                      )
                    })}
                  </div>
                )}
                </div>
              </label>
              <div flex flex-row>
              <h4 className="mb-2 text-xl font-semibold">Algorithm Type</h4>
              <div>
                <button
                  type="button"
                  className={dynamicStyle(activeAlgorithm === 'BFS')}
                  onClick={() => handleAlgorithmClick('BFS')}
                  data-twe-ripple-init
                  data-twe-ripple-color="light"
                >
                  BFS
                </button>
                <button
                  type="button"
                  className={dynamicStyle(activeAlgorithm === 'IDS')}
                  onClick={() => handleAlgorithmClick('IDS')}
                  data-twe-ripple-init
                  data-twe-ripple-color="light"
                >
                  IDS
                </button>
              </div>
              <h4 className="mb-2 text-xl font-semibold mt-3">Solution Type</h4>
              <div>
                <button
                  type="button"
                  className={dynamicStyle(activeSolution === 'Single Solution')}
                  onClick={() => handleSolutionClick('Single Solution')}
                  data-twe-ripple-init
                  data-twe-ripple-color="light"
                >
                  Single Solution
                </button>
                <button
                  type="button"
                  className={dynamicStyle(activeSolution === 'Multi Solution')}
                  onClick={() => handleSolutionClick('Multi Solution')}
                  data-twe-ripple-init
                  data-twe-ripple-color="light"
                >
                  Multi Solution
                </button>
              </div>
              <button type="submit"
                className='mt-4 mx-4 mb-[50px] rounded border-2 border-neutral-50 px-7 pb-[8px] pt-[10px] text-sm font-bold uppercase leading-normal text-neutral-50 transition duration-150 ease-in-out hover:border-neutral-100 hover:bg-neutral-500 hover:bg-opacity-10 hover:text-neutral-100 focus:border-neutral-100 focus:text-neutral-100 focus:outline-none focus:ring-0 active:border-neutral-200 active:text-neutral-200 dark:hover:bg-neutral-100 dark:hover:bg-opacity-10'
                data-twe-ripple-init
                data-twe-ripple-color="light">
                  Submit!
              </button>
              </div>
            </form>
            {loading && (
              <div className="flex justify-center items-center mt-4">
                <BeatLoader color="#ffffff" loading={loading} css={override} size={15} />
                <p className="ml-2 text-white">Loading...</p>
              </div>
            )}
            {submitted && (
              <div>
                <p className="text-white mt-4">Found <strong>... paths</strong> from <strong>{awal}</strong> to <strong>{akhir}</strong> in ... seconds!</p>
                <h2 className='mt-8 text-2xl font-bold'> Result </h2>
                <div className="mt-4 w-full flex flex-col items-center justify-center">
                  <div className="w-[900px] h-[450px] bg-black font-inter rounded-[10px] border-2 border-white mr-2"
                    style={{ backgroundColor: 'rgba(255, 255, 255, 0)' }}>
                    <p className="text-white">Result content here...</p>
                  </div>
                  <h2 className='mt-8 mb-4 text-2xl font-bold'> Individual Paths </h2>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
