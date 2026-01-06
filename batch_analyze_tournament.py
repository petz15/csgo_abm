"""
Batch analysis script for tournament matchup folders.
Automatically analyzes all matchup folders and generates HTML reports.
"""

import os
import json
import argparse
from pathlib import Path
from datetime import datetime
from concurrent.futures import ProcessPoolExecutor, as_completed
import sys
import papermill as pm

def find_matchup_folders(tournament_folder):
    """Find all matchup_XXX folders in the tournament directory."""
    matchup_folders = []
    tournament_path = Path(tournament_folder)
    
    if not tournament_path.exists():
        print(f"Error: Tournament folder '{tournament_folder}' not found!")
        return []
    
    for item in sorted(tournament_path.iterdir()):
        if item.is_dir() and item.name.startswith('matchup_'):
            csv_file = item / "all_games_minimal.csv"
            json_file = item / "simulation_summary.json"
            
            if csv_file.exists() and json_file.exists():
                matchup_folders.append(item)
            else:
                print(f"⚠️  Skipping {item.name} - missing CSV or JSON file")
    
    return matchup_folders

def analyze_matchup(matchup_folder, notebook_path):
    """
    Analyze a single matchup by executing the notebook with papermill.
    The notebook handles its own export and filename generation.
    Returns (success, matchup_name, error_msg)
    """
    matchup_name = matchup_folder.name
    
    try:
        # Paths
        csv_path = str((matchup_folder / "all_games_minimal.csv").resolve())
        json_path = str((matchup_folder / "simulation_summary.json").resolve())
        folder_path = str(matchup_folder.resolve())
        
        # Execute notebook using papermill
        # Parameters are injected directly into the notebook
        pm.execute_notebook(
            str(notebook_path),
            None,  # Don't save output notebook
            parameters={
                'FOLDER_PATH': folder_path,
                'CSV_FILE_PATH': csv_path,
                'CSV_INFO_FILE_PATH': json_path,
            },
            cwd=str(matchup_folder.parent),
            progress_bar=False,
            request_save_on_cell_execute=False,
            kernel_name='python3'
        )
        
        return (True, matchup_name, None)
            
    except pm.PapermillExecutionError as e:
        # Extract the relevant error information
        error_msg = str(e)
        if len(error_msg) > 500:
            error_msg = error_msg[-500:]
        return (False, matchup_name, error_msg)
    except Exception as e:
        return (False, matchup_name, str(e))

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(
        description='Batch analyze tournament matchup folders and generate HTML reports.',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''
Examples:
  python batch_analyze_tournament.py results_20260105_123550
  python batch_analyze_tournament.py results_20260105_123550 --workers 8
  python batch_analyze_tournament.py results_20260105_123550 --notebook analysis_notebook.ipynb
''')
    
    parser.add_argument('tournament_folder', 
                        help='Path to the tournament results folder containing matchup_XXX subdirectories')
    parser.add_argument('--notebook', '-n',
                        default='analysis_notebook_v2.ipynb',
                        help='Path to the Jupyter notebook to execute (default: analysis_notebook_v2.ipynb)')
    parser.add_argument('--workers', '-w',
                        type=int,
                        default=4,
                        help='Number of parallel workers (default: 4)')
    
    args = parser.parse_args()
    
    # Configuration from arguments
    TOURNAMENT_FOLDER = args.tournament_folder
    NOTEBOOK_PATH = Path(args.notebook).resolve()  # Convert to absolute path
    MAX_WORKERS = args.workers
    
    # Validate inputs
    if not Path(TOURNAMENT_FOLDER).exists():
        print(f"❌ Error: Tournament folder '{TOURNAMENT_FOLDER}' does not exist!")
        return 1
    
    if not NOTEBOOK_PATH.exists():
        print(f"❌ Error: Notebook '{NOTEBOOK_PATH}' does not exist!")
        return 1
    
    print("=" * 80)
    print("TOURNAMENT BATCH ANALYSIS (Parallel)")
    print("=" * 80)
    print(f"Tournament Folder: {TOURNAMENT_FOLDER}")
    print(f"Notebook: {NOTEBOOK_PATH}")
    print(f"Parallel Workers: {MAX_WORKERS}")
    print()
    
    # Find all matchup folders
    matchup_folders = find_matchup_folders(TOURNAMENT_FOLDER)
    
    if not matchup_folders:
        print("❌ No valid matchup folders found!")
        return 1
    
    print(f"✓ Found {len(matchup_folders)} matchup folders to analyze\n")
    
    # Track results
    successful = []
    failed = []
    start_time = datetime.now()
    
    # Process in parallel
    with ProcessPoolExecutor(max_workers=MAX_WORKERS) as executor:
        # Submit all jobs
        futures = {
            executor.submit(analyze_matchup, folder, NOTEBOOK_PATH): folder 
            for folder in matchup_folders
        }
        
        # Process as they complete
        for i, future in enumerate(as_completed(futures), 1):
            folder = futures[future]
            try:
                success, name, error = future.result()
                
                if success:
                    successful.append(name)
                    print(f"✓ [{i}/{len(matchup_folders)}] {name} - SUCCESS")
                else:
                    failed.append((name, error))
                    print(f"✗ [{i}/{len(matchup_folders)}] {name} - FAILED")
                    if error:
                        print(f"    Error: {error[:100]}")
            except Exception as e:
                failed.append((folder.name, str(e)))
                print(f"✗ [{i}/{len(matchup_folders)}] {folder.name} - EXCEPTION: {e}")
    
    # Summary
    elapsed = datetime.now() - start_time
    print("\n" + "=" * 80)
    print("BATCH ANALYSIS COMPLETE")
    print("=" * 80)
    print(f"Total Time: {elapsed}")
    print(f"Successful: {len(successful)}/{len(matchup_folders)}")
    print(f"Failed: {len(failed)}/{len(matchup_folders)}")
    
    if successful:
        print(f"\n✓ Successfully analyzed {len(successful)} matchups")
    
    if failed:
        print(f"\n✗ Failed matchups:")
        for name, error in failed[:10]:  # Show first 10 failures
            print(f"  - {name}: {error[:100]}")
        if len(failed) > 10:
            print(f"  ... and {len(failed) - 10} more")
    
    return 0 if not failed else 2  # Exit code 2 if some analyses failed



if __name__ == "__main__":
    sys.exit(main())
