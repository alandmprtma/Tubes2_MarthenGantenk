import Image from 'next/image';
import Link from 'next/link';
import icon from '../public/Lemanspedia-removebg.png';

const Navbar = () => {
    return (
    <nav className="sticky top-[0px] bg-[#333] py-[15px] z-20 mx-12 flex text-white items-center justify-around w-full"
    style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}>
        <Image src={icon} alt="Lemanspedia Logo" width={70} height={70} className='ml-[25px]' />
        <ul className='flex flex-row gap-x-[170px]'>
            <li>
                <Link href="/" className='hover:text-[#ADD8E6] cursor-pointer transition-none mx-8'>
                    <span className='font-bold text-xl'>Home</span>
                </Link>
            </li>
            <li>
                <Link href="/wiki-race" className='hover:text-[#ADD8E6] cursor-pointer transition-none mx-4'>
                    <span className='font-bold text-xl'>Wiki Race</span>
                </Link>
            </li>
            <li>
                <Link href="/about-us" className='hover:text-[#ADD8E6] cursor-pointer transition-none mx-8'>
                    <span className='font-bold text-xl'>About Us</span> 
                </Link>
            </li>
        </ul>
        <div />
    </nav>
    );
}

export default Navbar;