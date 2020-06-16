import json
import requests
import pymysql
from flask import Flask, render_template, request, session, flash

app = Flask(__name__)
#Change below if using different db provider - This is dockerised db.
db = pymysql.connect(
    host = '0.0.0.0',
    port = 3306,
    user = 'user',
    passwd = '1234',
    db = 'qse'
)
app.secret_key='Keyword'

@app.route('/', methods=['GET','POST'])
def show_search():
    #Enter search and press button
    return render_template('index.html')

@app.route('/results', methods=['GET','POST'])
def show_results(): 
    #Get keyword
    keyword = request.form['txtSearch']
    #Database logic
    cursor = db.cursor()
    #Setting up syntax
    searchParam = "%"+keyword+"%"
    search_query = "SELECT * FROM Pages WHERE Pages.Content LIKE %s"
    cursor.execute(search_query, searchParam)
    results=cursor.fetchall()
    ad_query = search_query = "SELECT Advert FROM Adverts WHERE Adverts.Keyword LIKE %s"
    cursor.execute(ad_query, searchParam)
    adResults=cursor.fetchall()
    
    return render_template('results.html', results=results, adverts=adResults)

if __name__ == "__main__":
    app.run(host = "0.0.0.0")