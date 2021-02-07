package main;

import java.util.Random;
import java.lang.*;

class Tournament extends Thread {
    int[] participants;
    volatile WrapInt winner;
    int begin;
    int end;

    Tournament(int[] monks, WrapInt winner, int begin, int end) {
        participants = monks;
        this.winner = winner;
        this.begin = begin;
        this.end = end;
    }

    void work() {
        if (end - begin == 1) {
            if (participants[begin] > participants[end]) {
                participants[begin] -= participants[end];
                winner.set(begin);
                return;
            } else {
                participants[end] -= participants[begin];
                winner.set(end);
                return;
            }
        }

        WrapInt leftWInner = new WrapInt(0);
        Tournament left = new Tournament(participants, leftWInner, begin, begin + (end - begin) / 2);
        left.start();

        WrapInt rightWinner = new WrapInt(0);
        Tournament right = new Tournament(participants, rightWinner, begin + (end - begin) / 2 + 1, end);
        right.work();

        try {
            left.join();
        } catch (InterruptedException err) {
            System.err.println(err);
        }

        if (participants[leftWInner.get()] > participants[rightWinner.get()]) {
            participants[leftWInner.get()] -= participants[rightWinner.get()];
            winner.set(leftWInner.get());
            return;
        } else {
            participants[rightWinner.get()] -= participants[leftWInner.get()];
            winner.set(rightWinner.get());
            return;
        }
    }

    public void start() {
        work();
    }
}

public class Main {

    static int[] generateMonks(int size, int winner) {
        int[] monks = new int[size];

        Random ran = new Random();

        for (int i = 0; i < size; i++) {
            monks[i] = ran.nextInt(100);
        }

        monks[winner] = 100_000_000;

        return monks;
    }
    public static void main(String[] args) {
        WrapInt winner = new WrapInt(0);
        int size = Integer.parseInt(args[0]);
        int winnerPos = Integer.parseInt(args[1]);

        Tournament t = new Tournament(generateMonks(size, winnerPos), winner, 0, size-1);
        t.work();

        System.out.println(winner.get());
    }
}

class WrapInt {
    Integer val = 0;

    public Integer get() {
        return val;
    }

    public void set(Integer val) {
        this.val = val;
    }

    public WrapInt(Integer val) {
        this.val = val;
    }
}