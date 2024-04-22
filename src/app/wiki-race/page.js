"use client"
import React, { useState } from 'react';
import Navbar from '../../components/navbar.js';
import { css } from "@emotion/react";
import { BeatLoader } from "react-spinners";
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

  const handleChangeAwal = (event) => {
    setAwal(event.target.value);
    setInputsFilled(event.target.value !== '' || akhir !== '');
  };

  const handleChangeAkhir = (event) => {
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

  return (
    <section
      className="rounded-lg bg-cover bg-no-repeat p-12 text-center relative object-cover"
      style={{backgroundImage: "url('/bg-website.png')", width: "100vw", height: "100vh"}}>
        {/* <Particles
          params={{
            "particles": {
              "number": {
                "value": 200,
                "density": {
                  "enable": true,
                  "value_area": 800
                }
              },
              "color": {
                "value": "#ffffff"
              },
              "shape": {
                "type": "circle",
                "stroke": {
                  "width": 0,
                  "color": "#000000"
                },
                "polygon": {
                  "nb_sides": 5
                }
              },
              "opacity": {
                "value": 0.5,
                "random": true,
                "anim": {
                  "enable": false,
                  "speed": 1,
                  "opacity_min": 0.1,
                  "sync": false
                }
              },
              "size": {
                "value": 3,
                "random": true,
                "anim": {
                  "enable": false,
                  "speed": 40,
                  "size_min": 0.1,
                  "sync": false
                }
              },
              "line_linked": {
                "enable": false,
                "distance": 150,
                "color": "#ffffff",
                "opacity": 0.4,
                "width": 1
              },
              "move": {
                "enable": true,
                "speed": 6,
                "direction": "none",
                "random": false,
                "straight": false,
                "out_mode": "out",
                "bounce": false,
                "attract": {
                  "enable": false,
                  "rotateX": 600,
                  "rotateY": 1200
                }
              }
            },
            "interactivity": {
              "detect_on": "canvas",
              "events": {
                "onhover": {
                  "enable": true,
                  "mode": "bubble"
                },
                "onclick": {
                  "enable": true,
                  "mode": "repulse"
                },
                "resize": true
              },
              "modes": {
                "grab": {
                  "distance": 400,
                  "line_linked": {
                    "opacity": 1
                  }
                },
                "bubble": {
                  "distance": 250,
                  "size": 0,
                  "duration": 2,
                  "opacity": 0,
                  "speed": 3
                },
                "repulse": {
                  "distance": 400,
                  "duration": 0.4
                },
                "push": {
                  "particles_nb": 4
                },
                "remove": {
                  "particles_nb": 2
                }
              }
            },
            "retina_detect": true
          }}
        /> */}
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
                <input
                  className='w-[500px] h-[60px] font-inter rounded-[10px] border-2 border-white mr-2'
                  style={{ backgroundColor: 'rgba(255, 255, 255, 0)' }}
                  type="text"
                  value={awal}
                  onChange={handleChangeAwal}
                  placeholder="Masukkan alamat awal"
                />
                {inputsFilled ? (
                  <img src="arrow-right-arrow-left-solid.svg" alt="to" className="w-[25px] mx-3 mt-4 mb-4" />
                ) : (
                  <h2 className="mx-3 mt-4 mb-4 text-2xl font-semibold">to</h2>
                )}
                <input
                  className='w-[500px] h-[60px] font-inter rounded-[10px] border-2 border-white mr-2'
                  style={{ backgroundColor: 'rgba(255, 255, 255, 0)' }}
                  type="text"
                  value={akhir}
                  onChange={handleChangeAkhir}
                  placeholder="Masukkan alamat akhir"
                />
              </label>
              <h4 className="mb-2 text-xl font-semibold">Algorithm Type</h4>
              <div flex flex-row>
                <div>
                  <button
                    type="button"
                    className="mx-4 rounded border-2 border-neutral-50 px-7 pb-[8px] pt-[10px] text-sm font-medium uppercase leading-normal text-neutral-50 transition duration-150 ease-in-out hover:border-neutral-100 hover:bg-neutral-500 hover:bg-opacity-10 hover:text-neutral-100 focus:border-neutral-100 focus:text-neutral-100 focus:outline-none focus:ring-0 active:border-neutral-200 active:text-neutral-200 dark:hover:bg-neutral-100 dark:hover:bg-opacity-10"
                    data-twe-ripple-init
                    data-twe-ripple-color="light">
                    BFS
                  </button>
                  <button
                    type="button"
                    className="mx-4 rounded border-2 border-neutral-50 px-7 pb-[8px] pt-[10px] text-sm font-medium uppercase leading-normal text-neutral-50 transition duration-150 ease-in-out hover:border-neutral-100 hover:bg-neutral-500 hover:bg-opacity-10 hover:text-neutral-100 focus:border-neutral-100 focus:text-neutral-100 focus:outline-none focus:ring-0 active:border-neutral-200 active:text-neutral-200 dark:hover:bg-neutral-100 dark:hover:bg-opacity-10"
                    data-twe-ripple-init
                    data-twe-ripple-color="light">
                    IDS
                  </button>
                </div>
              
                <button type="submit"
                className='mt-4 mx-4 rounded border-2 border-neutral-50 px-7 pb-[8px] pt-[10px] text-sm font-bold uppercase leading-normal text-neutral-50 transition duration-150 ease-in-out hover:border-neutral-100 hover:bg-neutral-500 hover:bg-opacity-10 hover:text-neutral-100 focus:border-neutral-100 focus:text-neutral-100 focus:outline-none focus:ring-0 active:border-neutral-200 active:text-neutral-200 dark:hover:bg-neutral-100 dark:hover:bg-opacity-10'
                data-twe-ripple-init
                data-twe-ripple-color="light">
                  GO!Lang
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
