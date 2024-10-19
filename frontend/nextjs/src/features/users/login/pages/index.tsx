"user client"

import { ChevronRight, ChevronDown } from 'lucide-react';
import React, { useState } from "react"
import { AxiosError } from "axios"
// import { useRouter } from "next/navigation"
import { getHttpClient } from "@/commons/setups/axios"
import { ErrorResponse } from '@/commons/helpers/response';

export function LoginPage() {
    // const router = useRouter()
    const [emailValue, setEmailOnChange] = useState("")
    const [passwordValue, setPasswordOnChange] = useState("")
    const [message, setMessage] = useState("")
    const [emailErrorMessage, setEmailErrorMessage] = useState("")
    const [passwordErrorMessage, setPasswordErrorMessage] = useState("")
    const [messageErrorMessage, setMessageErrorMessage] = useState("")
    const [pending, setPending] = useState(false)

    const setEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
        setEmailOnChange(e.target.value)
        setEmailErrorMessage("")
    }
    const setPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPasswordOnChange(e.target.value)
        setPasswordErrorMessage("")
    }

    const Login = async () => {
        // event: React.MouseEvent<HTMLButtonElement>
        // event.preventDefault();
        setPending(true)
        setMessage("")
        setMessageErrorMessage("")
        try {        
            const httpClient = getHttpClient()
            const requesstBody = {
                email: emailValue,
                password: passwordValue
            }
            const response = await httpClient.post("/api/v1/users/login", requesstBody, {
                headers: {
                    "Accept": "application/json"
                    // "X-REQUEST-ID": "requestId"
                }
            })
            setMessage(response.data.data.message)
            // console.log(response);
            
            // router.push("/")
        } catch(error: unknown) {
            if (error instanceof AxiosError) {
                error.response?.data.errors.forEach((element: ErrorResponse) => {
                    if (element.field === "email") {
                        setEmailErrorMessage(element.message)
                    } else if (element.field === "password") {
                        setPasswordErrorMessage(element.message)
                    } else if (element.field === "message") {
                        setMessageErrorMessage(element.message)
                    }
                });
            } else {
                console.log("error:", error);
            }
        }
        setPending(false)
    }

    return (
        <div className="flex justify-center min-h-screen">
            <div>
                <div className="w-full my-3 text-center">
                    ECOMMERCE
                </div>
                <div className="w-[350px] border rounded-lg p-5">
                    <div className="text-3xl mb-3">Sign in</div>
                    {message && <p className="text-center text-green-800 text-[0.80rem] mb-[15px]">{message}</p>}
                    {messageErrorMessage && <p className="text-center text-red-800 text-[0.80rem] mb-[15px]">{messageErrorMessage}</p>}
                    <label htmlFor="email" className="text-xs font-black mb-1">Email or mobile phone number</label>
                    <input type="text" id="email" name="email" 
                        disabled={pending} 
                        value={emailValue}
                        onChange={(e) => setEmail(e)}
                        className="w-full px-2 py-1 focus:outline-none focus:ring-2 text-sm border rounded-md border-gray-500"/>
                    {emailErrorMessage && <p className="text-red-300 text-[0.60rem]">{emailErrorMessage}</p>}
                    <label htmlFor="password" className="text-xs font-black mb-1">Password</label>
                    <input type="password" id="password" name="password" 
                        disabled={pending} 
                        value={passwordValue}
                        onChange={(e) => setPassword(e)}
                        className="w-full px-2 py-1 focus:outline-none focus:ring-2 text-sm border rounded-md border-gray-500"/>
                    {passwordErrorMessage && <p className="text-red-300 text-[0.60rem]">{passwordErrorMessage}</p>}
                    <button type="button" 
                        disabled={pending} 
                        onClick={Login}
                        className="bg-yellow-400 text-sm text-black font-standar py-2 px-4 rounded-md mt-3 w-full shadow
                            hover:bg-yellow-500 
                            active:bg-yellow-600 
                            disabled:bg-yellow-200 disabled:cursor-not-allowed">
                        Login
                    </button>
                    <span className="text-xs font-medium">By continuing, you agree to Ecommerce&apos;s <a href="#" className="text-blue-500 underline">Conditions of Use</a> and <a href="#" className="text-blue-500 underline">Privacy Notice</a>.</span>
                    <div className="mt-5 border">
                        <div className=""><ChevronRight className="w-3 h-3 invisible border"/> <ChevronDown className="w-3 h-3 visible border"/></div>
                    </div>
                    <div className="w-full border-t mt-3"></div>
                    <div className="text-xs font-black mt-3">Buying for work?</div>
                    <a href="#" className="text-xs text-blue-500 hover:underline">Shop on Ecommerce Business</a>
                </div>
                <div className="flex items-center justify-between w-full mt-5">
                    <div className="border-t w-[calc(100%-70%)]"></div>
                    <div className="text-xs flex-1 text-center text-gray-500">New to Ecommerce?</div>
                    <div className="border-t w-[calc(100%-70%)]"></div>
                </div>
                <button type="button" 
                    className="bg-white text-black text-sm font-standar border border-gray-300 px-4 py-2 rounded w-full mt-3
                        hover:bg-gray-100 
                        active:bg-gray-200 
                        disabled:bg-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed">
                    Create your account
                </button>
            </div>
        </div>
    )
}