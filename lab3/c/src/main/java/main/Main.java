package main;

import java.util.Random;
import java.util.concurrent.Semaphore;
import java.lang.*;


class Dealer extends Thread {
    public static volatile int[] board = new int[2];
    public static volatile Semaphore boardLock = new Semaphore(0);
    public static volatile Semaphore smokerLock = new Semaphore(0);
    public static volatile Semaphore midLock = new Semaphore(0);

    public void run() {
        int index = 0;

        while (true) {
            try {
                board[0] = index;
                index = (index + 1) % 3;
                board[1] = index;

                boardLock.release();

                smokerLock.acquire();

                boardLock.acquire();

                midLock.release();
                
            } catch (Exception err) {
                System.err.println(err);
            }
        }
    }
}

class Smoker extends Thread {
    int hasIngredient;

    public static volatile Semaphore nSmokerLock = new Semaphore(1);

    Smoker(int p) {
        hasIngredient = p;
    }
    
    public void run() {
        while (true) {
            try {
                nSmokerLock.acquire();
                Dealer.boardLock.acquire();

                if (Dealer.board[0] + Dealer.board[1] + hasIngredient == 3) {

                    System.out.println("Smoking " + Thread.currentThread().getName());
                    Thread.sleep(1000);

                    Dealer.smokerLock.release();

                    Dealer.boardLock.release();

                    Dealer.midLock.acquire();
                    
                    nSmokerLock.release();

                    continue;
                }

                Dealer.boardLock.release();
                nSmokerLock.release();
                
            } catch (Exception err) {
                System.err.println(err);
            }
        }
    }
}

public class Main {
    public static void main(String[] args) {
        Smoker s0 = new Smoker(0);
        Smoker s1 = new Smoker(1);
        Smoker s2 = new Smoker(2);

        Dealer d = new Dealer();
        d.start();

        s0.start();
        s1.start();
        s2.start();
    }
}