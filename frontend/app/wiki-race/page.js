"use client"
import React, { useState } from 'react';
import Navbar from '../../components/navbar.js';
import { css } from "@emotion/react";
import { BeatLoader } from "react-spinners";
import axios from 'axios';
import PathBox from '../../components/PathBox.js'
import Graph from '../../components/Graph.js'
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
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [openAwal, setOpenAwal] = useState(false);
  const [openAkhir, setOpenAkhir] = useState(false);
  // State untuk autocomplete
  const [resultAwal, setResultAwal]= useState([])
  const [resultAkhir, setResultAkhir]= useState([])
  const [activeAlgorithm, setActiveAlgorithm] = useState(''); // default kosong
  const [activeSolution, setActiveSolution] = useState(''); // default kosong
  const [errorMessage, setErrorMessage] = useState(null)

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
    setSubmitted(false);
    setInputsFilled(event.target.value !== '' || akhir !== '');
  };

  const handleChangeAkhir = (event) => {
    setOpenAkhir(true);
    handleQueryAkhir();
    setAkhir(event.target.value);
    setSubmitted(false);
    setInputsFilled(awal !== '' || event.target.value !== '');
  };

  const handleSubmit = async (event) => {
    setOpenAwal(false)
    setOpenAkhir(false)
    setSubmitted(false)
    event.preventDefault();
    
    if (awal =='' || akhir == ''){
      setErrorMessage("Please complete the start and the target.");
      setLoading(false)
      await delay(1500);
      setErrorMessage(null);
      return;
    }
    else if (activeAlgorithm == '' && activeSolution == ''){
      setErrorMessage("Please choose the algorithm and solution.");
      setLoading(false)
      await delay(1500);
      setErrorMessage(null);
      return;
    }
    else if (activeAlgorithm == ''){
      setErrorMessage("Please choose the algorithm.");
      setLoading(false)
      await delay(1500);
      setErrorMessage(null);
      return;
    }
    else if (activeSolution == ''){
      setErrorMessage("Please choose the solution.");
      setLoading(false);
      await delay(1500);
      setErrorMessage(null);
      return;
    }
    setLoading(true);
    try {
      const response = await axios.post('http://localhost:8080/search', {
        start: awal,
        target: akhir,
        algorithm: activeAlgorithm,
        solution: activeSolution
      });
      setLoading(false);
      setSubmitted(true);
      setResults(response.data);
      setErrorMessage(null)
    } catch (error) {
      console.error(error);
      setLoading(false);
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
    <div className="flex flex-col w-full items-center">
      <Navbar/>
      <img src="/Lemanspedia_Slogan-removebg.png" className='h-[200px] object-cover'/>
      <h2 className="mt-[25px] mb-4 text-2xl font-semibold">Find the shortest paths from</h2>
      <div className="text-white justify-center">
        <form onSubmit={handleSubmit}>
          <label className=' mb-[20px] flex flex-row items-center justify-center'>
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
                    <div className='border-[1px] w-[475px] border-white flex flex-row items-center p-2 mb-2 rounded-lg'
                    style={{ backgroundColor: 'rgba(255, 255, 255, 0.1)' }}
                      key={index}
                      onClick={() => {
                        setAwal(item.title)
                        setOpenAwal((prev) =>!prev)
                      }}
                    >
                      {item.thumbnail && (
                        <img src={item.thumbnail} alt={item.title} className='w-10 h-10 mr-2 rounded-full object-cover'/>
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
              <h2 className="mx-3 mt-4 mb-4 text-2xl font-semibold h-[60px] flex items-center justify-center">to</h2>
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
                    <div className='border-[1px] w-[475px] border-white flex flex-row items-center p-2 mb-2 rounded-lg'
                    style={{ backgroundColor: 'rgba(255, 255, 255, 0.1)' }}
                      key={index}
                      onClick={() => {
                        setAkhir(item.title)
                        setOpenAkhir((prev) =>!prev)
                      }}
                    >
                      {item.thumbnail && (
                        <img src={item.thumbnail} alt={item.title} className='w-10 h-10 mr-2 rounded-full object-cover'/>
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
          <div className='flex flex-col items-center'>
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
          {!loading && (<button type="submit"
            className='mt-4 mx-4 mb-[15px] rounded border-2 border-neutral-50 px-7 pb-[8px] pt-[10px] text-sm font-bold uppercase leading-normal text-neutral-50 transition duration-150 ease-in-out hover:border-neutral-100 hover:bg-neutral-500 hover:bg-opacity-10 hover:text-neutral-100 focus:border-neutral-100 focus:text-neutral-100 focus:outline-none focus:ring-0 active:border-neutral-200 active:text-neutral-200 dark:hover:bg-neutral-100 dark:hover:bg-opacity-10'
            data-twe-ripple-init
            data-twe-ripple-color="light">
              Submit!
          </button>
          )}
          </div>
        </form>
        {errorMessage && (
          <div className='text-white text-center mb-4'>{errorMessage}</div>
        )}
        {loading && (
          <div className="flex justify-center items-center mt-[25px] mb-[50px]">
            <BeatLoader color="#ffffff" loading={loading} css={override} size={15} />
            <p className="ml-2 text-white">Loading...</p>
          </div>
        )}
        {submitted && results && (
          <div className='flex flex-col items-center'>
              {results.numberOfPaths === 0 ? (
                  <p className="text-white text-center mt-4 text-xl">No path found from <strong>{awal}</strong> to <strong>{akhir}</strong></p>
              ) : (
                  <>
                      <div className='w-[50%]'>
                          <p className="text-white text-center mt-4 text-xl">Found <strong>{results.numberOfPaths} paths</strong> from <strong>{awal}</strong> to <strong>{akhir}</strong> in <strong>{results.elapsedTime} seconds</strong>!</p>
                          <p className="text-white text-center mt-4 text-xl">Articles Checked: <strong>{results.articlesChecked}</strong></p>
                          <p className="text-white text-center mt-4 text-xl">Articles Traversed: <strong>{results.articlesTraversed}</strong></p>
                      </div>
                      <div className='w-[1140px] bg-white h-[2px] mt-2'/>
                      <h2 className='mt-5 text-2xl font-bold'> Connecting Graphs </h2>
                      <div className='w-[900px] h-[450px] font-inter rounded-[10px] border-2 border-white mr-2 overflow-hidden'>
                          <div className='translate-x-[-200x] translate-y-[100px] z-[-10px]'>
                              <Graph path={results.paths}/>
                          </div>
                          <div className='flex translate-y-[-760px] w-[150px] z-[10px] h-fit rounded-[10px] border-2 border-white mt-2 ml-2'>
                              <p>Drag to pan. Scroll to zoom.</p>
                          </div>
                      </div>
                      <div className='w-[1140px] bg-white h-[2px] mt-4'/>
                      <h2 className='mt-5 text-2xl font-bold'> Individual Paths </h2>
                      <div className=" w-full flex flex-col items-center justify-center">
                          <PathBox path={results.paths} />
                      </div>
                  </>
              )}
          </div>
        )}
      </div>
    </div>
  );
}
