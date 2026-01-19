from langchain_ollama import ChatOllama
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_core.runnables import RunnablePassthrough
from src.config import settings
from src.infra.qdrant_client import get_vector_store


class RAGengine :
    def __init__(self):
        self.vector_store = get_vector_store()
        self.llm = ChatOllama(
            base_url=settings.OLLAMA_BASE_URL,
            model=settings.OLLAMA_MODEL
        )
        #we will be searching for the top 3 docs
        self.retriever = self.vector_store.as_retriever(search_kwargs={"k":3})
        
    def _format_docs(self , docs):
        #all the data from the docs that we extracted is called here to give to LLM as context
        return "\n\n".join(f"[Source: {d.metadata.get('title','unkown')}] {d.page_content}" for d in docs)
    
    def query(self,question: str) -> str:
        #now here itself I can apply multiple types of Query translations and modify how I want my RAG to behave
        template = """
        Answer the question based only on the following context:
        {context}
        
        Question : 
        {question}
        """
        prompt = ChatPromptTemplate.from_template(template)
        
        #The LCEL chan 
        
        chain = (
            {"context" : self.retriever | self._format_docs, "question" : RunnablePassthrough()} #Preapre the inputs 
            | prompt #format the prompt
            | self.llm #ask the LLM
            | StrOutputParser() #Clean the output
        )
        
        return chain.invoke(question)