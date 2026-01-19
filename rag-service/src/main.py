import argparse
import logging
import sys
import os


sys.path.append(os.getcwd())

from src.core.engine import RAGengine

logging.basicConfig(level=logging.INFO , format='%(asctime)s - %(levelname)s - %(message)s ')
logger = logging.getLogger(__name__)

def main():
    parser = argparse.ArgumentParser(description="Poduction RAG Query Service")
    parser.add_argument("query" , type=str , help="the question to ask")
    args = parser.parse_args()
    
    try:
        logger.info(f"Initializing RAG engine....")
        engine = RAGengine()
        logger.info(f"Analyzing Query: {args.query}")
        answer = engine.query(args.query)
        
        print("\n" + "="*50)
        print("🤖 STRATEGY REPORT")
        print("="*50)
        print(answer)
        print("="*50 + "\n")
    except Exception as e:
        logger.error(f"System Error : {e}" , exc_info=True)
        sys.exit(1)
        

if __name__ == "__main__":
    main()